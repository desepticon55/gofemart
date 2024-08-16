package withdrawal

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type withdrawalService interface {
	FindAllWithdrawals(ctx context.Context) ([]model.Withdrawal, error)
}
