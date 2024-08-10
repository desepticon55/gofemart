package balance

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type balanceRepository interface {
	FindBalance(ctx context.Context, userName string) (common.Balance, error)

	FindBalanceStats(ctx context.Context, userName string) (common.BalanceStats, error)

	Withdraw(ctx context.Context, balance common.Balance, sum float64, orderNumber string) error
}
