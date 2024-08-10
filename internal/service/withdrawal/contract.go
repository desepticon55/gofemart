package withdrawal

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type withdrawalRepository interface {
	FindAllWithdrawals(ctx context.Context, userName string) ([]common.Withdrawal, error)
}
