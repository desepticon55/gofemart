package balance

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type balanceService interface {
	FindBalanceStats(ctx context.Context) (model.BalanceStats, error)

	Withdraw(ctx context.Context, orderNumber string, sum float64) error
}
