package withdrawal

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type withdrawalService interface {
	FindAllWithdrawals(ctx context.Context) ([]common.Withdrawal, error)
}
