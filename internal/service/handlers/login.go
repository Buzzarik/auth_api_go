package handlers

import (
	my_hash "auth/internal/lib/hash"
	"auth/internal/lib/jwt"
	"auth/internal/service"
	"fmt"
	"log/slog"
	"net/http"
)

func LoginHandlers(app *service.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.LoginHandlers";

		var input struct{
			PhoneNumber 	string 	`json:"phone_number"`
			Password 		string 	`json:"password"`
			ID_API			int64	`json:"id_api"`
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


		if input.Password == "" || input.PhoneNumber == "" {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Phone and password are required");
			app.Log.Info("Error no content request");
			return;
		}

		//запрашиваем данные из storageUser + хешируем пароль
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


		if user == nil || !my_hash.ComparePassword(user.HashPassword, input.Password){
			app.ErrorResponse(w, http.StatusBadRequest,
				"Wrong login or password");
			app.Log.Info("Wrong login or password",
					slog.String("phone_number", input.PhoneNumber));
			return;
		}

		token, err := jwt.NewToken(user, app);
		token.IdAPI = input.ID_API;

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed create Token");
			app.Log.Debug("Failed create Token",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		//TODO: token save in DB
		err = app.StorageToken.SetToken(token);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed set token from storageToken");
			app.Log.Debug("Failed set token from storageToken",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		err = app.WriteJSON(w, http.StatusCreated, 
			service.Envelope{
				"success": true,
				"token": token,
				"message": "Token created successfully",
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

		app.Log.Info("Token created or exists",
			slog.String("token", fmt.Sprintf("%v",token)));
	}
}