package handlers

import (
	"auth/internal/service"
	"net/http"
)

func RecoverPanic(app *service.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.ErrorResponse(w, http.StatusInternalServerError, "Failed to recover")
			}
		}()

		next.ServeHTTP(w, r)
	})
}