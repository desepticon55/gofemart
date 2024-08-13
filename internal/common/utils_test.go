package common

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsValidOrderNumber(t *testing.T) {
	t.Run("valid order number", func(t *testing.T) {
		validOrderNumber := "1234567812345670"
		result := IsValidOrderNumber(validOrderNumber)
		assert.True(t, result)
	})

	t.Run("invalid order number", func(t *testing.T) {
		invalidOrderNumber := "1234567812345671"
		result := IsValidOrderNumber(invalidOrderNumber)
		assert.False(t, result)
	})

	t.Run("order number with non-numeric characters", func(t *testing.T) {
		invalidOrderNumber := "1234A67812345670"
		result := IsValidOrderNumber(invalidOrderNumber)
		assert.False(t, result)
	})

	t.Run("empty order number", func(t *testing.T) {
		emptyOrderNumber := ""
		result := IsValidOrderNumber(emptyOrderNumber)
		assert.False(t, result)
	})
}

func TestCreateJWTToken(t *testing.T) {
	t.Run("should create a valid JWT token", func(t *testing.T) {
		username := "testUser"
		tokenString, err := CreateJWTToken(username)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(*Claims)
		assert.True(t, ok)
		assert.Equal(t, username, claims.Username)

		expectedExpiration := time.Now().Add(5 * time.Minute).Truncate(time.Second)
		actualExpiration := claims.ExpiresAt.Time.Truncate(time.Second)
		assert.WithinDuration(t, expectedExpiration, actualExpiration, time.Second)
	})
}
