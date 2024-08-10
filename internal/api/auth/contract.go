package auth

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
)

type userService interface {
	CreateUser(ctx context.Context, user common.User) error

	FindUser(ctx context.Context, userName string) (common.User, error)
}
