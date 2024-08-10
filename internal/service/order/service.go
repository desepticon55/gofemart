package order

import (
	"context"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"go.uber.org/zap"
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
		return common.ErrOrderNumberIsNotFilled
	}

	if !common.IsValidOrderNumber(orderNumber) {
		s.logger.Error("Order number is not valid", zap.String("orderNumber", orderNumber))
		return common.ErrOrderNumberIsNotValid
	}

	exist, err := s.orderRepository.ExistOrder(ctx, orderNumber)
	if err != nil {
		s.logger.Error("Error during check exist order", zap.String("orderNumber", orderNumber), zap.Error(err))
		return err
	}

	currentUserName := fmt.Sprintf("%v", ctx.Value(common.UserNameContextKey))
	if !exist {
		err := s.orderRepository.CreateOrder(ctx, common.Order{
			OrderNumber:    orderNumber,
			Username:       currentUserName,
			CreateDate:     time.Now(),
			LastModifyDate: time.Now(),
			Status:         common.NewOrderStatus,
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
			return common.ErrOrderNumberHasUploadedCurrentUser
		} else {
			return common.ErrOrderNumberHasUploadedOtherUser
		}
	}

	return nil
}

func (s *OrderService) FindAllOrders(ctx context.Context) ([]common.Order, error) {
	currentUserName := fmt.Sprintf("%v", ctx.Value(common.UserNameContextKey))
	orders, err := s.orderRepository.FindAllOrders(ctx, currentUserName)
	if err != nil {
		s.logger.Error("Error during find orders", zap.String("userName", currentUserName), zap.Error(err))
		return nil, err
	}

	if len(orders) == 0 {
		return nil, common.ErrOrdersWasNotFound
	}

	return orders, nil
}
