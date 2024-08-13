package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	pool, cleanup := internal.InitPostgresIntegrationTest(t, ctx, logger)

	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Fatalf("failed to cleanup test database: %s", err)
		}
	})

	orderRepository := NewOrderRepository(pool, logger)

	order := common.Order{
		OrderNumber:    "12345678903",
		CreateDate:     time.Now(),
		LastModifyDate: time.Now(),
		Status:         "NEW",
		Username:       "testUser",
		Accrual:        40.,
		KeyHash:        10,
		KeyHashModule:  0,
		Version:        0,
	}

	t.Run("ExistOrder", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := orderRepository.CreateOrder(ctx, order)
		assert.NoError(t, err)

		result, err := orderRepository.ExistOrder(ctx, "12345678903")
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("CreateOrder", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := orderRepository.CreateOrder(ctx, order)
		assert.NoError(t, err)
	})

	t.Run("FindOrder", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := orderRepository.CreateOrder(ctx, order)
		assert.NoError(t, err)

		result, err := orderRepository.FindOrder(ctx, "12345678903")
		assert.NoError(t, err)
		assert.Equal(t, order.OrderNumber, result.OrderNumber)
		assert.Equal(t, order.Status, result.Status)
		assert.Equal(t, order.Username, result.Username)
		assert.Equal(t, order.Accrual, result.Accrual)
		assert.Equal(t, order.KeyHash, result.KeyHash)
		assert.Equal(t, order.KeyHashModule, result.KeyHashModule)
		assert.Equal(t, order.Version, result.Version)
	})

	t.Run("FindAllOrders", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})

		err := orderRepository.CreateOrder(ctx, order)
		assert.NoError(t, err)

		result, err := orderRepository.FindAllOrders(ctx, "testUser")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, order.OrderNumber, result[0].OrderNumber)
		assert.Equal(t, order.Status, result[0].Status)
		assert.Equal(t, order.Username, result[0].Username)
		assert.Equal(t, order.Accrual, result[0].Accrual)
		assert.Equal(t, order.KeyHash, result[0].KeyHash)
		assert.Equal(t, order.KeyHashModule, result[0].KeyHashModule)
		assert.Equal(t, order.Version, result[0].Version)
	})

	t.Run("ChangeOrderStatus", func(t *testing.T) {
		t.Cleanup(func() {
			if err := internal.ClearTables(ctx, pool); err != nil {
				t.Fatalf("failed to clear tables: %s", err)
			}
		})
		if _, err := pool.Exec(ctx, `INSERT INTO gofemart.balance (username, balance, opt_lock) VALUES ($1, $2, $3)`, "testUser", 100, 0); err != nil {
			t.Fatalf("failed to insert balance: %v", err)
		}

		err := orderRepository.CreateOrder(ctx, order)
		assert.NoError(t, err)

		err = orderRepository.ChangeOrderStatus(ctx, order, "PROCESSED", 555.0)
		assert.NoError(t, err)

		result, err := orderRepository.FindOrder(ctx, "12345678903")
		assert.NoError(t, err)
		assert.Equal(t, order.OrderNumber, result.OrderNumber)
		assert.Equal(t, "PROCESSED", result.Status)
		assert.Equal(t, order.Username, result.Username)
		assert.Equal(t, 555.0, result.Accrual)
		assert.Equal(t, order.KeyHash, result.KeyHash)
		assert.Equal(t, order.KeyHashModule, result.KeyHashModule)
		assert.Equal(t, int64(1), result.Version)

		var balance float64
		err = pool.QueryRow(ctx, `SELECT balance FROM gofemart.balance WHERE username = $1`, "testUser").Scan(&balance)
		assert.NoError(t, err)
		assert.Equal(t, 655., balance)
	})
}
