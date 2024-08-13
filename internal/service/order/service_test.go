package order

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"math"
	"testing"
	"time"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) ExistOrder(ctx context.Context, orderNumber string) (bool, error) {
	args := m.Called(ctx, orderNumber)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order common.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindOrder(ctx context.Context, orderNumber string) (common.Order, error) {
	args := m.Called(ctx, orderNumber)
	return args.Get(0).(common.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAllOrders(ctx context.Context, userName string) ([]common.Order, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).([]common.Order), args.Error(1)
}

func TestOrderService_UploadOrder(t *testing.T) {

	t.Run("should return error if order number is empty", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		err := service.UploadOrder(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberIsNotFilled, err)
	})

	t.Run("should return error if order number is invalid", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		invalidOrderNumber := "invalid"
		err := service.UploadOrder(ctx, invalidOrderNumber)
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberIsNotValid, err)
	})

	t.Run("should return error because check exist order return error", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		mockRepo.On("ExistOrder", ctx, orderNumber).Return(false, errors.New("db error"))

		err := service.UploadOrder(ctx, orderNumber)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should create order", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		currentTime := time.Now()
		keyHash := int64(math.Abs(float64(common.HashCode(orderNumber))))

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(false, nil)
		mockRepo.On("CreateOrder", ctx, common.Order{
			OrderNumber:    orderNumber,
			Username:       "testUser",
			CreateDate:     currentTime,
			LastModifyDate: currentTime,
			Status:         common.NewOrderStatus,
			KeyHash:        keyHash,
			KeyHashModule:  keyHash % int64(common.Module),
		}).Return(nil)

		err := service.UploadOrder(ctx, orderNumber)
		assert.NoError(t, err)
	})

	t.Run("should return error if order exists to current user", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		existingOrder := common.Order{
			OrderNumber: orderNumber,
			Username:    "testUser",
		}

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(true, nil)
		mockRepo.On("FindOrder", ctx, orderNumber).Return(existingOrder, nil)

		err := service.UploadOrder(ctx, orderNumber)
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberHasUploadedCurrentUser, err)
	})

	t.Run("should return error if order exists to another user", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		existingOrder := common.Order{
			OrderNumber: orderNumber,
			Username:    "otherUser",
		}

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(true, nil)
		mockRepo.On("FindOrder", ctx, orderNumber).Return(existingOrder, nil)

		err := service.UploadOrder(ctx, orderNumber)
		assert.Error(t, err)
		assert.Equal(t, common.ErrOrderNumberHasUploadedOtherUser, err)
	})
}

func TestOrderService_FindAllOrders(t *testing.T) {
	t.Run("should return error if there is an error during find orders", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		mockRepo.On("FindAllOrders", ctx, "testUser").Return([]common.Order{}, errors.New("db error"))

		orders, err := service.FindAllOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should return error if no orders are found", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		mockRepo.On("FindAllOrders", ctx, "testUser").Return([]common.Order{}, nil)

		orders, err := service.FindAllOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, common.ErrOrdersWasNotFound, err)
	})

	t.Run("should return orders successfully", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), common.UserNameContextKey, "testUser")

		service := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orders := []common.Order{
			{OrderNumber: "12345"},
			{OrderNumber: "67890"},
		}
		mockRepo.On("FindAllOrders", ctx, "testUser").Return(orders, nil)

		resultOrders, err := service.FindAllOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, orders, resultOrders)
	})
}
