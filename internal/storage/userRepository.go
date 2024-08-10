package storage

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type UserRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *UserRepository) ExistUser(ctx context.Context, userName string) (bool, error) {
	var exist bool
	query := "select exists(select 1 from gofemart.user where username = $1)"
	err := r.pool.QueryRow(ctx, query, userName).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, userName string, password string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		r.logger.Error("Error during open transaction", zap.Error(err))
		return err
	}

	query := "insert into gofemart.user(username, password) values ($1, $2)"
	_, err = tx.Exec(ctx, query, userName, password)
	if err != nil {
		r.logger.Error("Error during create user", zap.String("userName", userName), zap.Error(err))
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

func (r *UserRepository) FindUser(ctx context.Context, userName string) (common.User, error) {
	var user common.User
	query := "select username, password from gofemart.user where username = $1"
	err := r.pool.QueryRow(ctx, query, userName).Scan(&user.Username, &user.Password)
	if err != nil {
		return common.User{}, err
	}

	return user, nil
}
