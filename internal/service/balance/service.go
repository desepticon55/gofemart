package balance

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"go.uber.org/zap"
)

type BalanceService struct {
	logger            *zap.Logger
	balanceRepository balanceRepository
}

func NewBalanceService(l *zap.Logger, r balanceRepository) *BalanceService {
	return &BalanceService{logger: l, balanceRepository: r}
}

func (s *BalanceService) FindBalanceStats(ctx context.Context) (common.BalanceStats, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(common.UserNameContextKey))
	balance, err := s.balanceRepository.FindBalanceStats(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during fetch balance", zap.String("userName", currentUserName), zap.Error(err))
		return common.BalanceStats{}, err
	}
	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, orderNumber string, sum float64) error {
	if orderNumber == "" || sum == 0 {
		return common.ErrOrderNumberOrSumIsNotFilled
	}

	currentUserName := fmt.Sprintf("%v", ctx.Value(common.UserNameContextKey))
	if !common.IsValidOrderNumber(orderNumber) {
		return common.ErrOrderNumberIsNotValid
	}

	balance, err := s.balanceRepository.FindBalance(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during fetch balance", zap.String("userName", currentUserName), zap.Error(err))
		return err
	}

	if balance.Balance < sum {
		return common.ErrUserBalanceLessThanSumToWithdraw
	}

	err = s.balanceRepository.Withdraw(ctx, balance, sum, orderNumber)
	if err != nil {
		s.logger.Error("Error during withdraw", zap.String("userName", currentUserName), zap.String("orderNumber", orderNumber), zap.Error(err))
		return err
	}

	return nil
}
