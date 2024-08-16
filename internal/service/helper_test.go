package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
