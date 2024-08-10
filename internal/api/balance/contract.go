package balance

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type balanceService interface {
	FindBalanceStats(ctx context.Context) (common.BalanceStats, error)

	Withdraw(ctx context.Context, orderNumber string, sum float64) error
}
