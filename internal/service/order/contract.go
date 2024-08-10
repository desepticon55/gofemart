package order

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type orderRepository interface {
	ExistOrder(ctx context.Context, orderNumber string) (bool, error)

	FindOrder(ctx context.Context, orderNumber string) (common.Order, error)

	CreateOrder(ctx context.Context, order common.Order) error

	FindAllOrders(ctx context.Context, userName string) ([]common.Order, error)
}
