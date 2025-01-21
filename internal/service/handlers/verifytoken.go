package handlers

import (
	"auth/internal/lib/jwt"
	"auth/internal/service"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func VerifyTokenHandlers(app *service.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handlers.VerifyTokenHandlers";

		var input struct {
			HashToken 	string 	`json:"token"`
			ID_API 		int 	`json:"id_api"`	
			IdUser		int		`json:"id_user"`
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

		if input.HashToken == "" {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Token is required");
			app.Log.Info("Error no content request");
			return;
		}

		//удаляем невалидные токены
		err = app.StorageToken.DeleteToken(time.Now());
		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed get user from StorageToken");
			app.Log.Debug("Failed get user from StorageToken",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}


		token, err := jwt.DecodeToken(input.HashToken, app);

		if err != nil {
			app.ErrorResponse(w, http.StatusNotFound,
				"Invalid token or expiry");
			app.Log.Info("token is not Valid");
			app.Log.Debug("token is not Valid",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		dbToken, err := app.StorageToken.SelectOneToken(token.IdUser, input.ID_API);
		
		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Failed get token from StorageToken");
			app.Log.Debug("Failed get token from StorageToken",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		if dbToken == nil {
			app.ErrorResponse(w, http.StatusNotFound,
				"Invalid token or expiry");
			app.Log.Info("Token is not exists");
			return;
		}


		//сравниваем хеш токена и id апи, откуда пришел запрос с БД
		if !jwt.VerifyToken(input.HashToken, input.ID_API, dbToken) {
			app.ErrorResponse(w, http.StatusBadRequest,
				"Invalid token or expiry");
			app.Log.Info("token is not Valid. Hash not EQ");
			return;
		}


		err = app.WriteJSON(w, http.StatusOK, 
			service.Envelope{
				"success": true,
				"message": "Token valid",
			}, nil);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Error write JSON in response");
			app.Log.Debug("Error write JSON in response",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		app.Log.Info("Token valid",
			slog.String("token", fmt.Sprintf("%v", dbToken)));

	}
}