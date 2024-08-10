package common

import (
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type ContextKey string

const (
	UserNameContextKey ContextKey = "userName"
)

var JwtKey = []byte("hard_coded_jwt_secret_key")

func CreateJWTToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func IsValidOrderNumber(orderNumber string) bool {
	sum := 0
	needDouble := false

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}
		if needDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		needDouble = !needDouble
	}
	return sum%10 == 0
}
