package withdrawal

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

type MockWithdrawalRepository struct {
	mock.Mock
}

func (m *MockWithdrawalRepository) FindAllWithdrawals(ctx context.Context, userName string) ([]model.Withdrawal, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).([]model.Withdrawal), args.Error(1)
}

func TestWithdrawalService_TestFindAllWithdrawals(t *testing.T) {
	ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")
	logger := zaptest.NewLogger(t)

	t.Run("should return found withdrawals", func(t *testing.T) {
		mockRepo := new(MockWithdrawalRepository)
		service := &WithdrawalService{
			withdrawalRepository: mockRepo,
			logger:               logger,
		}

		expectedWithdrawals := []model.Withdrawal{
			{ID: "1", Username: "testUser", OrderNumber: "123456", Sum: 100.0},
			{ID: "2", Username: "testUser", OrderNumber: "123456", Sum: 200.0},
		}
		mockRepo.On("FindAllWithdrawals", ctx, "testUser").Return(expectedWithdrawals, nil)

		withdrawals, err := service.FindAllWithdrawals(ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedWithdrawals, withdrawals)

		mockRepo.AssertCalled(t, "FindAllWithdrawals", ctx, "testUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when database return error", func(t *testing.T) {
		mockRepo := new(MockWithdrawalRepository)
		service := &WithdrawalService{
			withdrawalRepository: mockRepo,
			logger:               logger,
		}

		expectedError := errors.New("database error")
		mockRepo.On("FindAllWithdrawals", ctx, "testUser").Return([]model.Withdrawal{}, expectedError)

		withdrawals, err := service.FindAllWithdrawals(ctx)
		assert.Error(t, err)
		assert.Empty(t, withdrawals)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertCalled(t, "FindAllWithdrawals", ctx, "testUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when withdrawals was not found", func(t *testing.T) {
		mockRepo := new(MockWithdrawalRepository)
		service := &WithdrawalService{
			withdrawalRepository: mockRepo,
			logger:               logger,
		}

		mockRepo.On("FindAllWithdrawals", ctx, "testUser").Return([]model.Withdrawal{}, nil)

		withdrawals, err := service.FindAllWithdrawals(ctx)
		assert.Error(t, err)
		assert.Empty(t, withdrawals)
		assert.Equal(t, model.ErrWithdrawalsWasNotFound, err)

		mockRepo.AssertCalled(t, "FindAllWithdrawals", ctx, "testUser")
		mockRepo.AssertExpectations(t)
	})
}
