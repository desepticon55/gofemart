package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type TransactionFunc func(tx pgx.Tx) error

func transactional(ctx context.Context, logger *zap.Logger, pool *pgxpool.Pool, fn TransactionFunc) (err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error during open transaction: %w", err)
	}

	defer func() {
		if err != nil {
			logger.Error("Error during rollback transaction", zap.Error(err))
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)
	return err
}
