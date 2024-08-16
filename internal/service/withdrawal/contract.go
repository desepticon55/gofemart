package withdrawal

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type withdrawalRepository interface {
	FindAllWithdrawals(ctx context.Context, userName string) ([]model.Withdrawal, error)
}
