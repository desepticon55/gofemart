package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type WithdrawalRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewWithdrawalRepository(pool *pgxpool.Pool, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *OrderRepository) FindAllWithdrawals(ctx context.Context, userName string) ([]common.Withdrawal, error) {
	query := "select id, order_number, username, sum, create_date from gofemart.withdrawal where username = $1 order by create_date"
	rows, err := r.pool.Query(ctx, query, userName)
	if err != nil {
		r.logger.Error("Error during execute query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var withdrawals []common.Withdrawal
	for rows.Next() {
		var withdrawal common.Withdrawal
		if err := rows.Scan(&withdrawal.ID, &withdrawal.OrderNumber, &withdrawal.Username, &withdrawal.Sum, &withdrawal.CreateDate); err != nil {
			r.logger.Error("Error during scan row", zap.Error(err))
			continue
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
