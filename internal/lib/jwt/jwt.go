package jwt

import (
	"auth/internal/models"
	"auth/internal/service"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken creates new JWT token for given user and app.
func NewToken(user *models.User, app *service.Application) (*models.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256);
	expiry := time.Now().Add(app.Cnf.Server.TokenTTL);
	claims := token.Claims.(jwt.MapClaims);
	claims["phone_number"] = user.PhoneNumber;
	claims["name"] = user.Name;
	claims["created_at"] = user.CreatedAt;
	claims["exp"] = expiry;


	tokenString, err := token.SignedString([]byte(app.Cnf.Server.Secret))
	if err != nil {
		return nil, err;
	}

	return &models.Token{
		Hash: tokenString,
		Expiry: expiry,
		IdAPI: app.Cnf.Server.IdAPI,
		PhoneNumber: user.PhoneNumber,
		Name: user.Name,
	}, nil;
}