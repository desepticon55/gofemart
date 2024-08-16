package auth

import (
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateJWTToken(t *testing.T) {
	t.Run("should create a valid JWT token", func(t *testing.T) {
		username := "testUser"
		tokenString, err := createJWTToken(username)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(*model.Claims)
		assert.True(t, ok)
		assert.Equal(t, username, claims.Username)

		expectedExpiration := time.Now().Add(5 * time.Minute).Truncate(time.Second)
		actualExpiration := claims.ExpiresAt.Time.Truncate(time.Second)
		assert.WithinDuration(t, expectedExpiration, actualExpiration, time.Second)
	})
}
