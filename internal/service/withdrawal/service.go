package withdrawal

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"go.uber.org/zap"
)

type WithdrawalService struct {
	logger               *zap.Logger
	withdrawalRepository withdrawalRepository
}

func NewWithdrawalService(l *zap.Logger, r withdrawalRepository) *WithdrawalService {
	return &WithdrawalService{logger: l, withdrawalRepository: r}
}

func (s *WithdrawalService) FindAllWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(service.UserNameContextKey))
	withdrawals, err := s.withdrawalRepository.FindAllWithdrawals(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during find withdrawals", zap.String("userName", currentUserName), zap.Error(err))
		return []model.Withdrawal{}, err
	}

	if len(withdrawals) == 0 {
		return []model.Withdrawal{}, model.ErrWithdrawalsWasNotFound
	}

	return withdrawals, nil
}
