package order

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"go.uber.org/zap"
	"math"
	"time"
)

type OrderService struct {
	logger          *zap.Logger
	orderRepository orderRepository
}

func NewOrderService(l *zap.Logger, r orderRepository) *OrderService {
	return &OrderService{logger: l, orderRepository: r}
}

func (s *OrderService) UploadOrder(ctx context.Context, orderNumber string) error {
	if orderNumber == "" {
		s.logger.Error("Order number is not filled")
		return model.ErrOrderNumberIsNotFilled
	}

	if !service.IsValidOrderNumber(orderNumber) {
		s.logger.Error("Order number is not valid", zap.String("orderNumber", orderNumber))
		return model.ErrOrderNumberIsNotValid
	}

	exist, err := s.orderRepository.ExistOrder(ctx, orderNumber)
	if err != nil {
		s.logger.Error("Error during check exist order", zap.String("orderNumber", orderNumber), zap.Error(err))
		return err
	}

	currentUserName := fmt.Sprintf("%v", ctx.Value(service.UserNameContextKey))
	if !exist {
		keyHash := int64(math.Abs(float64(service.HashCode(orderNumber))))
		err := s.orderRepository.CreateOrder(ctx, model.Order{
			OrderNumber:    orderNumber,
			Username:       currentUserName,
			CreateDate:     time.Now(),
			LastModifyDate: time.Now(),
			Status:         model.NewOrderStatus,
			KeyHash:        keyHash,
			KeyHashModule:  keyHash % int64(service.Module),
		})
		if err != nil {
			s.logger.Error("Error during create order", zap.String("orderNumber", orderNumber), zap.Error(err))
			return err
		}
	} else {
		order, err := s.orderRepository.FindOrder(ctx, orderNumber)
		if err != nil {
			s.logger.Error("Error during find order", zap.String("orderNumber", orderNumber), zap.Error(err))
			return err
		}
		if currentUserName == order.Username {
			return model.ErrOrderNumberHasUploadedCurrentUser
		} else {
			return model.ErrOrderNumberHasUploadedOtherUser
		}
	}

	return nil
}

func (s *OrderService) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(service.UserNameContextKey))
	orders, err := s.orderRepository.FindAllOrders(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during find orders", zap.String("userName", currentUserName), zap.Error(err))
		return nil, err
	}

	if len(orders) == 0 {
		return nil, model.ErrOrdersWasNotFound
	}

	return orders, nil
}
