package withdrawal

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"go.uber.org/zap"
)

type WithdrawalService struct {
	logger               *zap.Logger
	withdrawalRepository withdrawalRepository
}

func NewWithdrawalService(l *zap.Logger, r withdrawalRepository) *WithdrawalService {
	return &WithdrawalService{logger: l, withdrawalRepository: r}
}

func (s *WithdrawalService) FindAllWithdrawals(ctx context.Context) ([]common.Withdrawal, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(common.UserNameContextKey))
	withdrawals, err := s.withdrawalRepository.FindAllWithdrawals(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during find withdrawals", zap.String("userName", currentUserName), zap.Error(err))
		return []common.Withdrawal{}, err
	}

	if len(withdrawals) == 0 {
		return []common.Withdrawal{}, common.ErrWithdrawalsWasNotFound
	}

	return withdrawals, nil
}
