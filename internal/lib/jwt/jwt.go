package jwt

import (
	"auth/internal/models"
	"auth/internal/service"
	"fmt"
	_"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
    PhoneNumber string `json:"phone_number"`
    Name        string `json:"name"`
    CreatedAt   time.Time `json:"created_at"`
    Exp         time.Time `json:"exp"`
    IdUser      int64 `json:"id_user"`
	jwt.RegisteredClaims
}

// NewToken creates new JWT token for given user and app.
func NewToken(user *models.User, app *service.Application) (*models.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256);
	expiry := time.Now().Add(app.Cnf.Server.TokenTTL);
	claims := token.Claims.(jwt.MapClaims);
	claims["phone_number"] = user.PhoneNumber;
	claims["name"] = user.Name;
	claims["created_at"] = user.CreatedAt;
	claims["exp"] = expiry;
	claims["id_user"] = user.ID;


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
		IdUser: user.ID,
	}, nil;
}

func DecodeToken(tokenString string, app *service.Application) (*models.Token, error) {
    mySigningKey := []byte(app.Cnf.Server.Secret) // Секретный ключ
    claims := &MyCustomClaims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return mySigningKey, nil
    })
    if err != nil {
        return nil, fmt.Errorf("error parsing token: %w", err)
    }

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

    return &models.Token{
        Hash:        tokenString,
        Expiry:     claims.Exp,
        PhoneNumber: claims.PhoneNumber,
        Name:        claims.Name,
        IdUser:      claims.IdUser,
    }, nil
}

func VerifyToken(hash string, id_api int64, token *models.Token) bool {
	return token.Hash == hash && id_api == token.IdAPI;
}