package order

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type orderService interface {
	UploadOrder(ctx context.Context, orderNumber string) error

	FindAllOrders(ctx context.Context) ([]model.Order, error)
}
