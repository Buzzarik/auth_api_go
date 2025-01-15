package handlers

import (
	my_hash "auth/internal/lib/hash"
	"auth/internal/lib/otp"
	"auth/internal/service"
	"log/slog"
	"net/http"
)

func SignupUserHandler(app *service.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.SignupUserHandler";

		var input struct {
			Name string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			Password string `json:"password"`
		};

		err := app.ReadJSON(w, r, &input);

		if err != nil {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Invalid request payload");
			app.Log.Error("Error reading JSON in request");
			app.Log.Debug("Error reading JSON in request",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		if input.Name == "" || input.PhoneNumber == "" || input.Password == "" {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Name and phone and password number are required");
			app.Log.Info("Error no content request");
			return;
		}

		user, err := app.StorageUser.GetByPhoneNumber(input.PhoneNumber);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed get user from StorageUser");
			app.Log.Debug("Failed get user from StorageUser",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		if user != nil {
			app.ErrorResponse(w, http.StatusConflict,
				"User already exists with the given phone number");
			app.Log.Info("User already exists",
					slog.String("phone_number", input.PhoneNumber));
			return;
		}

		//генерируем одноразовый (5 мин) пароль
		otp, err := otp.GenerateOTP();
		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed generation OTP");
			app.Log.Debug("Failed generation OTP",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		//записываем в cache
		//TODO: сделать хеширование пароля
		hash, err := my_hash.HashPassword(input.Password);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed hashing password");
			app.Log.Debug("Failed hashing password",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		userData := map[string]string{
			"name": input.Name,
			"otp":  otp,
			"hash_password": hash,
		};

		err = app.Cache.SetUser(input.PhoneNumber, userData);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed set user from Cache");
			app.Log.Debug("Failed set user from Cache",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		//TODO: здесь можно отправить ОТР на телефон, для подтвеждения и тд
		app.Log.Debug("Generated OTP", 
		slog.String("phone_number", input.PhoneNumber),
		slog.String("OTP", otp));
		err = app.WriteJSON(w, http.StatusOK, service.Envelope{
			"success": true,
			"message": "OTP sent successfully"},
		nil);

		if err != nil {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Server internal error");
			app.Log.Error("Error write JSON in response");
			app.Log.Debug("Error write JSON in response",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		app.Log.Info("OTP sent successfully. Response succes");
	}
}