package balance

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"go.uber.org/zap"
)

type BalanceService struct {
	logger            *zap.Logger
	balanceRepository balanceRepository
}

func NewBalanceService(l *zap.Logger, r balanceRepository) *BalanceService {
	return &BalanceService{logger: l, balanceRepository: r}
}

func (s *BalanceService) FindBalanceStats(ctx context.Context) (model.BalanceStats, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(service.UserNameContextKey))
	balance, err := s.balanceRepository.FindBalanceStats(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during fetch balance", zap.String("userName", currentUserName), zap.Error(err))
		return model.BalanceStats{}, err
	}
	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, orderNumber string, sum float64) error {
	if orderNumber == "" || sum == 0 {
		return model.ErrOrderNumberOrSumIsNotFilled
	}

	currentUserName := fmt.Sprintf("%v", ctx.Value(service.UserNameContextKey))
	if !service.IsValidOrderNumber(orderNumber) {
		return model.ErrOrderNumberIsNotValid
	}

	balance, err := s.balanceRepository.FindBalance(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during fetch balance", zap.String("userName", currentUserName), zap.Error(err))
		return err
	}

	if balance.Balance < sum {
		return model.ErrUserBalanceLessThanSumToWithdraw
	}

	err = s.balanceRepository.Withdraw(ctx, balance, sum, orderNumber)
	if err != nil {
		s.logger.Error("Error during withdraw", zap.String("userName", currentUserName), zap.String("orderNumber", orderNumber), zap.Error(err))
		return err
	}

	return nil
}
