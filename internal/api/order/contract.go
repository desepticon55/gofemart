package order

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type orderService interface {
	UploadOrder(ctx context.Context, orderNumber string) error

	FindAllOrders(ctx context.Context) ([]common.Order, error)
}
