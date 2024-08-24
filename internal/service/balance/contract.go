package balance

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type balanceRepository interface {
	FindBalance(ctx context.Context, userName string) (model.Balance, error)

	FindBalanceStats(ctx context.Context, userName string) (model.BalanceStats, error)

	Withdraw(ctx context.Context, balance model.Balance, sum float64, orderNumber string) error
}
