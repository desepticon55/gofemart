package order

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type orderRepository interface {
	ExistOrder(ctx context.Context, orderNumber string) (bool, error)

	FindOrder(ctx context.Context, orderNumber string) (model.Order, error)

	CreateOrder(ctx context.Context, order model.Order) error

	FindAllOrders(ctx context.Context, userName string) ([]model.Order, error)
}
