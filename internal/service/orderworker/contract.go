package orderworker

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type orderRepository interface {
	FindOrdersToProcess(ctx context.Context, from int, to int) ([]model.Order, error)

	ChangeOrderStatus(ctx context.Context, order model.Order, status string, accrual float64) error
}
