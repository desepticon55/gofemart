package storage

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type OrderRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewOrderRepository(pool *pgxpool.Pool, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r OrderRepository) ExistOrder(ctx context.Context, orderNumber string) (bool, error) {
	var exist bool
	query := "select exists(select 1 from gofemart.order where order_number = $1)"
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r OrderRepository) FindOrder(ctx context.Context, orderNumber string) (common.Order, error) {
	var order common.Order
	query := "select order_number, username, create_date, last_modify_date, status, accrual from gofemart.order where order_number = $1"
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(&order.OrderNumber, &order.Username, &order.CreateDate, &order.LastModifyDate, &order.Status, &order.Accrual)
	if err != nil {
		return common.Order{}, err
	}

	return order, nil
}

func (r OrderRepository) CreateOrder(ctx context.Context, order common.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("Error during open transaction", zap.Error(err))
		return err
	}

	query := "insert into gofemart.order(order_number, username, create_date, last_modify_date, status, accrual) values ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(ctx, query, order.OrderNumber, order.Username, order.CreateDate, order.LastModifyDate, order.Status, order.Accrual)
	if err != nil {
		r.logger.Error("Error during create order", zap.String("orderNumber", order.OrderNumber), zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Error during commit transaction. Start rollback", zap.Error(err))
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			r.logger.Error("Error during rollback transaction", zap.Error(err))
		}
		return errors.Join(err, rollbackErr)
	}

	return nil
}

func (r OrderRepository) FindAllOrders(ctx context.Context, userName string) ([]common.Order, error) {
	query := "select order_number, username, create_date, last_modify_date, status, accrual from gofemart.order where username = $1 order by create_date"
	rows, err := r.pool.Query(ctx, query, userName)
	if err != nil {
		r.logger.Error("Error during execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []common.Order
	for rows.Next() {
		var order common.Order
		if err := rows.Scan(&order.OrderNumber, &order.Username, &order.CreateDate, &order.LastModifyDate, &order.Status, &order.Accrual); err != nil {
			r.logger.Error("Error during scan row", zap.Error(err))
			continue
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
