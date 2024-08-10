package orderworker

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type orderRepository interface {
	FindOrdersToProcess(ctx context.Context, from int, to int) ([]common.Order, error)

	ChangeOrderStatus(ctx context.Context, order common.Order, status string, accrual float64) error
}
