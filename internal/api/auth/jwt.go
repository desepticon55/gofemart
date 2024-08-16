package auth

import (
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var JwtKey = []byte("hard_coded_jwt_secret_key")

func createJWTToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &model.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
