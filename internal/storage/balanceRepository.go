package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"time"
)

type BalanceRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewBalanceRepository(pool *pgxpool.Pool, logger *zap.Logger) *BalanceRepository {
	return &BalanceRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *BalanceRepository) FindBalance(ctx context.Context, userName string) (common.Balance, error) {
	query := "select username, balance, opt_lock from gofemart.balance where username = $1"
	var balance common.Balance
	err := r.pool.QueryRow(ctx, query, userName).Scan(&balance.Username, &balance.Balance, &balance.Version)
	if err != nil {
		return common.Balance{}, err
	}

	return balance, nil
}

func (r *BalanceRepository) FindBalanceStats(ctx context.Context, userName string) (common.BalanceStats, error) {
	query := `
		select b.username, b.balance, coalesce(sum(w.sum), 0) 
		from gofemart.balance b
		left join gofemart.withdrawal w on b.username = w.username
		where b.username = $1
		group by b.username, b.balance
    `
	var balance common.BalanceStats
	err := r.pool.QueryRow(ctx, query, userName).Scan(&balance.Username, &balance.Balance, &balance.Withdrawn)
	if err != nil {
		return common.BalanceStats{}, err
	}

	return balance, nil
}

func (r *BalanceRepository) Withdraw(ctx context.Context, balance common.Balance, sum float64, orderNumber string) error {
	return common.Transactional(ctx, r.pool, func(tx pgx.Tx) error {
		query := "update gofemart.balance set balance = $1, opt_lock = $2 where username = $3 and opt_lock = $4"
		result, err := tx.Exec(ctx, query, balance.Balance-sum, balance.Version+1, balance.Username, balance.Version)
		if err != nil {
			r.logger.Error("Error during change balance", zap.String("userName", balance.Username), zap.Error(err))
			return err
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			r.logger.Error("User balance has changed in other transaction", zap.String("userName", balance.Username), zap.Error(err))
			return common.ErrUserBalanceHasChanged
		}
		withdrawID, err := uuid.NewRandom()
		if err != nil {
			r.logger.Error("Error during generate UUID", zap.Error(err))
			return err
		}

		withdrawQuery := "insert into gofemart.withdrawal(id, order_number, username, sum, create_date) values ($1, $2, $3, $4, $5)"
		_, err = tx.Exec(ctx, withdrawQuery, withdrawID, orderNumber, balance.Username, sum, time.Now())
		if err != nil {
			r.logger.Error("Error during create withdrawal", zap.String("orderNumber", orderNumber), zap.Error(err))
			return err
		}
		return nil
	})
}
