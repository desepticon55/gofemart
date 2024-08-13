package balance

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"
)

type MockBalanceRepository struct {
	mock.Mock
}

func (m *MockBalanceRepository) FindBalanceStats(ctx context.Context, userName string) (common.BalanceStats, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(common.BalanceStats), args.Error(1)
}

func (m *MockBalanceRepository) FindBalance(ctx context.Context, userName string) (common.Balance, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(common.Balance), args.Error(1)
}

func (m *MockBalanceRepository) Withdraw(ctx context.Context, balance common.Balance, sum float64, orderNumber string) error {
	args := m.Called(ctx, balance, sum, orderNumber)
	return args.Error(0)
}

func TestBalanceService_FindBalanceStats(t *testing.T) {
	t.Run("should return error if fetch balance return error", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		mockRepo.On("FindBalanceStats", ctx, "testUser").Return(common.BalanceStats{}, errors.New("db error"))

		stats, err := service.FindBalanceStats(ctx)
		assert.Error(t, err)
		assert.Equal(t, common.BalanceStats{}, stats)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should return balance stats successfully", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		expectedStats := common.BalanceStats{Username: "testUser", Balance: 1000, Withdrawn: 500}
		mockRepo.On("FindBalanceStats", ctx, "testUser").Return(expectedStats, nil)

		stats, err := service.FindBalanceStats(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
	})
}

func TestBalanceService_Withdraw(t *testing.T) {
	t.Run("should return error if order number is empty or sum is zero", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		err := service.Withdraw(ctx, "", 0)
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberOrSumIsNotFilled, err)
	})

	t.Run("should return error if order number is invalid", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		orderNumber := "invalid"

		err := service.Withdraw(ctx, orderNumber, 100)
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberIsNotValid, err)
	})

	t.Run("should return error if fetch balance return error", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		orderNumber := "12345678903"
		mockRepo.On("FindBalance", ctx, "testUser").Return(common.Balance{}, errors.New("db error"))

		err := service.Withdraw(ctx, orderNumber, 100)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should return error if balance is less than sum to withdraw", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		orderNumber := "12345678903"
		mockRepo.On("FindBalance", ctx, "testUser").Return(common.Balance{Username: "testUser", Balance: 50}, nil)

		err := service.Withdraw(ctx, orderNumber, 100)
		assert.Error(t, err)
		assert.Equal(t, common.ErrUserBalanceLessThanSumToWithdraw, err)
	})

	t.Run("should successfully withdraw", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockBalanceRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &BalanceService{
			logger:            logger,
			balanceRepository: mockRepo,
		}

		mockRepo.On("FindBalance", ctx, "testUser").Return(common.Balance{Username: "testUser", Balance: 200}, nil)
		mockRepo.On("Withdraw", ctx, common.Balance{Username: "testUser", Balance: 200}, 100.0, "12345678903").Return(nil)

		err := service.Withdraw(ctx, "12345678903", 100.0)
		assert.NoError(t, err)
	})
}
