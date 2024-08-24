package user

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
)

type userRepository interface {
	ExistUser(ctx context.Context, userName string) (bool, error)

	CreateUser(ctx context.Context, userName string, password string) error

	FindUser(ctx context.Context, userName string) (model.User, error)
}
