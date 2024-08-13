package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestBalanceRepository(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	pool, cleanup := internal.InitPostgresIntegrationTest(t, ctx, logger)
	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatalf("failed to cleanup test database: %s", err)
		}
	})

	balanceRepository := NewBalanceRepository(pool, logger)

	balance := common.Balance{
		Username: "testuser",
		Balance:  1000.0,
		Version:  1,
	}

	t.Run("FindBalance", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		if _, err := pool.Exec(ctx, `INSERT INTO gofemart.balance (username, balance, opt_lock) VALUES ($1, $2, $3)`,
			balance.Username, balance.Balance, balance.Version); err != nil {
			t.Fatalf("failed to insert balance: %v", err)
		}

		result, err := balanceRepository.FindBalance(ctx, "testuser")
		assert.NoError(t, err)
		assert.Equal(t, balance, result)
	})

	t.Run("Withdraw", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		if _, err := pool.Exec(ctx, `INSERT INTO gofemart.balance (username, balance, opt_lock) VALUES ($1, $2, $3)`,
			balance.Username, balance.Balance, balance.Version); err != nil {
			t.Fatalf("failed to insert balance: %v", err)
		}

		err := balanceRepository.Withdraw(ctx, balance, 200, "12345678903")
		assert.NoError(t, err)

		updatedBalance, err := balanceRepository.FindBalance(ctx, "testuser")
		assert.NoError(t, err)
		assert.Equal(t, 800., updatedBalance.Balance)

		var count int
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM gofemart.withdrawal WHERE order_number = $1`, "12345678903").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
