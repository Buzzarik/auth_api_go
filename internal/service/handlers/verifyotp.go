package handlers

import (
	"auth/internal/models"
	"auth/internal/service"
	"log/slog"
	"net/http"
)

const constraint = "CONSTRAINT";

func VerifyOTPHandler(app *service.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.VerifyOTP";

		var input struct {
			OTP string `json:"otp"`
			PhoneNumber string `json:"phone_number"`
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

		if input.OTP == "" || input.PhoneNumber == "" {
			app.ErrorResponse(w, http.StatusBadRequest,
				"OTP and phone are required");
			app.Log.Info("Error no content request");
			return;
		}

		//запрашиваем данные из cache
		userData, err := app.Cache.GetUser(input.PhoneNumber);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed get user from Cache");
			app.Log.Debug("Failed get user from Cache",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}
		if len(userData) == 0 { //вышло время временного пароля
			app.ErrorResponse(w, http.StatusNotFound,
				"expired OTP");
			app.Log.Info("expired OTP");
			return;
		}

		if input.OTP != userData["otp"] {
			app.ErrorResponse(w, http.StatusNotFound,
				"Invalid OTP");
			app.Log.Info("Invalid OTP");
			return;
		}

		user := &models.User{
			Name: userData["name"],
			HashPassword: userData["hash_password"],
			PhoneNumber: input.PhoneNumber,
		};

		err = app.StorageUser.SetUser(user);

		if err != nil && err.Error() == constraint {
			app.ErrorResponse(w, http.StatusBadRequest,
				"User with phone_number already exests");
			app.Log.Info("User with phone_number already exists");
			return;
		}

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed get user from storageUser");
			app.Log.Debug("Failed get user from storageUser",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		err = app.WriteJSON(w, http.StatusCreated, 
			service.Envelope{
				"success": true,
				"name":  user.Name,
				"phone": user.PhoneNumber,
				"message": "User registered successfully",
			}, nil);
		
		if err != nil {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Server internal error");
			app.Log.Error("Error write JSON in response");
			app.Log.Debug("Error write JSON in response",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		app.Log.Info("User register",
			slog.String("name", user.Name),
			slog.String("phone_number", user.PhoneNumber));
	}
}