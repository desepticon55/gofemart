package storage

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/jackc/pgx/v4"
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
	return transactional(ctx, r.logger, r.pool, func(tx pgx.Tx) error {
		query := "insert into gofemart.user(username, password) values ($1, $2)"
		_, err := tx.Exec(ctx, query, userName, password)
		if err != nil {
			r.logger.Error("Error during create user", zap.String("userName", userName), zap.Error(err))
			return err
		}

		balanceQuery := "insert into gofemart.balance(username, balance, opt_lock) values ($1, $2, $3)"
		_, err = tx.Exec(ctx, balanceQuery, userName, 0, 0)
		if err != nil {
			r.logger.Error("Error during create balance", zap.String("userName", userName), zap.Error(err))
			return err
		}
		return nil
	})
}

func (r *UserRepository) FindUser(ctx context.Context, userName string) (model.User, error) {
	var user model.User
	query := "select username, password from gofemart.user where username = $1"
	err := r.pool.QueryRow(ctx, query, userName).Scan(&user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
