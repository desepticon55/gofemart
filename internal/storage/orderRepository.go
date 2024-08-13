package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/jackc/pgx/v4"
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

func (r *OrderRepository) ExistOrder(ctx context.Context, orderNumber string) (bool, error) {
	var exist bool
	query := "select exists(select 1 from gofemart.order where order_number = $1)"
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r *OrderRepository) FindOrder(ctx context.Context, orderNumber string) (common.Order, error) {
	var order common.Order
	query := "select order_number, username, create_date, last_modify_date, status, accrual, opt_lock, key_hash, key_hash_module from gofemart.order where order_number = $1"
	err := r.pool.QueryRow(ctx, query, orderNumber).Scan(&order.OrderNumber, &order.Username, &order.CreateDate,
		&order.LastModifyDate, &order.Status, &order.Accrual, &order.Version, &order.KeyHash, &order.KeyHashModule)
	if err != nil {
		return common.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order common.Order) error {
	return common.Transactional(ctx, r.pool, func(tx pgx.Tx) error {
		query := `insert into gofemart.order(order_number, username, create_date, last_modify_date, status, accrual, key_hash, key_hash_module, opt_lock)
			      values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := tx.Exec(ctx, query, order.OrderNumber, order.Username, order.CreateDate, order.LastModifyDate,
			order.Status, order.Accrual, order.KeyHash, order.KeyHashModule, 0)
		if err != nil {
			r.logger.Error("Error during create order", zap.String("orderNumber", order.OrderNumber), zap.Error(err))
			return err
		}
		return nil
	})
}

func (r *OrderRepository) FindAllOrders(ctx context.Context, userName string) ([]common.Order, error) {
	query := `select order_number, username, create_date, last_modify_date, status, accrual, key_hash, key_hash_module, opt_lock
			  from gofemart.order where username = $1 order by create_date`
	rows, err := r.pool.Query(ctx, query, userName)
	if err != nil {
		r.logger.Error("Error during execute query", zap.Error(err))
		return []common.Order{}, err
	}
	defer rows.Close()

	var orders []common.Order
	for rows.Next() {
		var order common.Order
		if err := rows.Scan(&order.OrderNumber, &order.Username, &order.CreateDate, &order.LastModifyDate,
			&order.Status, &order.Accrual, &order.KeyHash, &order.KeyHashModule, &order.Version); err != nil {
			r.logger.Error("Error during scan row", zap.Error(err))
			continue
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return []common.Order{}, err
	}

	return orders, nil
}

func (r *OrderRepository) FindOrdersToProcess(ctx context.Context, from, to int) ([]common.Order, error) {
	query := `select order_number, username, create_date, last_modify_date, status, accrual, key_hash, key_hash_module, opt_lock 
              from gofemart.order
			  where status in ('NEW', 'PROCESSING') and key_hash_module >= $1 and key_hash_module < $2
			  limit 100
    `
	rows, err := r.pool.Query(ctx, query, from, to)
	if err != nil {
		r.logger.Error("Error during execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []common.Order
	for rows.Next() {
		var order common.Order
		if err := rows.Scan(&order.OrderNumber, &order.Username, &order.CreateDate, &order.LastModifyDate, &order.Status,
			&order.Accrual, &order.KeyHash, &order.KeyHashModule, &order.Version); err != nil {
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

func (r *OrderRepository) ChangeOrderStatus(ctx context.Context, order common.Order, status string, accrual float64) error {
	return common.Transactional(ctx, r.pool, func(tx pgx.Tx) error {
		if status == common.ProcessedOrderStatus {
			findBalanceQuery := "select username, balance, opt_lock from gofemart.balance where username = $1"
			var balance common.Balance
			err := r.pool.QueryRow(ctx, findBalanceQuery, order.Username).Scan(&balance.Username, &balance.Balance, &balance.Version)
			if err != nil {
				r.logger.Error("Error during find balance", zap.String("userName", balance.Username), zap.Error(err))
				return err
			}

			changeBalanceQuery := "update gofemart.balance set balance = $1, opt_lock = $2 where username = $3 and opt_lock = $4"
			result, err := tx.Exec(ctx, changeBalanceQuery, balance.Balance+accrual, balance.Version+1, balance.Username, balance.Version)
			if err != nil {
				r.logger.Error("Error during change balance", zap.String("userName", balance.Username), zap.Error(err))
				return err
			}

			rowsAffected := result.RowsAffected()
			if rowsAffected == 0 {
				r.logger.Error("User balance has changed in other transaction", zap.String("userName", balance.Username), zap.Error(err))
				return common.ErrUserBalanceHasChanged
			}
		}

		changeOrderQuery := "update gofemart.order set status = $1, accrual = $2, opt_lock = $3 where order_number = $4 and opt_lock = $5"
		result, err := tx.Exec(ctx, changeOrderQuery, status, accrual, order.Version+1, order.OrderNumber, order.Version)
		if err != nil {
			r.logger.Error("Error during change order", zap.String("orderNumber", order.OrderNumber), zap.Error(err))
			return err
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			r.logger.Error("Order has changed in other transaction", zap.String("orderNumber", order.OrderNumber), zap.Error(err))
			return common.ErrUserBalanceHasChanged
		}

		return nil
	})
}
