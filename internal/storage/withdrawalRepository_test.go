package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func TestWithdrawalRepository(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	pool, cleanup := internal.InitPostgresIntegrationTest(t, ctx, logger)

	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatalf("failed to cleanup test database: %s", err)
		}
	})

	withdrawalRepository := NewWithdrawalRepository(pool, logger)

	t.Run("FindAllWithdrawals", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		if _, err := pool.Exec(ctx, `INSERT INTO gofemart.withdrawal (id, order_number, username, sum, create_date) VALUES ($1, $2, $3, $4, $5)`,
			"c6c2a5b1-5c3b-4a70-a18b-7e1a7397c118", "12345678903", "testUser", 45., time.Now()); err != nil {
			t.Fatalf("failed to insert balance: %v", err)
		}

		result, err := withdrawalRepository.FindAllWithdrawals(ctx, "testUser")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, "c6c2a5b1-5c3b-4a70-a18b-7e1a7397c118", result[0].ID)
		assert.Equal(t, "12345678903", result[0].OrderNumber)
		assert.Equal(t, "testUser", result[0].Username)
		assert.Equal(t, 45., result[0].Sum)
	})
}
