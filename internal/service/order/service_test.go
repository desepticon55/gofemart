package order

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
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

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindOrder(ctx context.Context, orderNumber string) (model.Order, error) {
	args := m.Called(ctx, orderNumber)
	return args.Get(0).(model.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAllOrders(ctx context.Context, userName string) ([]model.Order, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).([]model.Order), args.Error(1)
}

func TestOrderService_UploadOrder(t *testing.T) {

	t.Run("should return error if order number is empty", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		err := orderService.UploadOrder(ctx, "")
		assert.Error(t, err)
		assert.Equal(t, model.ErrOrderNumberIsNotFilled, err)
	})

	t.Run("should return error if order number is invalid", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		invalidOrderNumber := "invalid"
		err := orderService.UploadOrder(ctx, invalidOrderNumber)
		assert.Error(t, err)
		assert.Equal(t, model.ErrOrderNumberIsNotValid, err)
	})

	t.Run("should return error because check exist order return error", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		mockRepo.On("ExistOrder", ctx, orderNumber).Return(false, errors.New("db error"))

		err := orderService.UploadOrder(ctx, orderNumber)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should create order", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		currentTime := time.Now()
		keyHash := int64(math.Abs(float64(service.HashCode(orderNumber))))

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(false, nil)
		mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("model.Order")).Return(nil).Run(func(args mock.Arguments) {
			order := args.Get(1).(model.Order)

			assert.Equal(t, orderNumber, order.OrderNumber)
			assert.Equal(t, "testUser", order.Username)
			assert.Equal(t, model.NewOrderStatus, order.Status)
			assert.Equal(t, keyHash, order.KeyHash)
			assert.Equal(t, keyHash%int64(service.Module), order.KeyHashModule)

			assert.WithinDuration(t, currentTime, order.CreateDate, time.Second)
			assert.WithinDuration(t, currentTime, order.LastModifyDate, time.Second)
		})

		err := orderService.UploadOrder(ctx, orderNumber)
		assert.NoError(t, err)
	})

	t.Run("should return error if order exists to current user", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		existingOrder := model.Order{
			OrderNumber: orderNumber,
			Username:    "testUser",
		}

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(true, nil)
		mockRepo.On("FindOrder", ctx, orderNumber).Return(existingOrder, nil)

		err := orderService.UploadOrder(ctx, orderNumber)
		assert.Error(t, err)
		assert.Equal(t, model.ErrOrderNumberHasUploadedCurrentUser, err)
	})

	t.Run("should return error if order exists to another user", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orderNumber := "12345678903"
		existingOrder := model.Order{
			OrderNumber: orderNumber,
			Username:    "otherUser",
		}

		mockRepo.On("ExistOrder", ctx, orderNumber).Return(true, nil)
		mockRepo.On("FindOrder", ctx, orderNumber).Return(existingOrder, nil)

		err := orderService.UploadOrder(ctx, orderNumber)
		assert.Error(t, err)
		assert.Equal(t, model.ErrOrderNumberHasUploadedOtherUser, err)
	})
}

func TestOrderService_FindAllOrders(t *testing.T) {
	t.Run("should return error if there is an error during find orders", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		mockRepo.On("FindAllOrders", ctx, "testUser").Return([]model.Order{}, errors.New("db error"))

		orders, err := orderService.FindAllOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("should return error if no orders are found", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		mockRepo.On("FindAllOrders", ctx, "testUser").Return([]model.Order{}, nil)

		orders, err := orderService.FindAllOrders(ctx)
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Equal(t, model.ErrOrdersWasNotFound, err)
	})

	t.Run("should return orders successfully", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		mockRepo := new(MockOrderRepository)
		ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")

		orderService := &OrderService{
			logger:          logger,
			orderRepository: mockRepo,
		}

		orders := []model.Order{
			{OrderNumber: "12345"},
			{OrderNumber: "67890"},
		}
		mockRepo.On("FindAllOrders", ctx, "testUser").Return(orders, nil)

		resultOrders, err := orderService.FindAllOrders(ctx)
		assert.NoError(t, err)
		assert.Equal(t, orders, resultOrders)
	})
}
