package user

import (
	"context"
	"github.com/desepticon55/gofemart/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger     *zap.Logger
	repository userRepository
}

func NewUserService(l *zap.Logger, r userRepository) *UserService {
	return &UserService{logger: l, repository: r}
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) error {
	if user.Username == "" || user.Password == "" {
		return model.ErrUserDataIsNotValid
	}

	exist, err := s.repository.ExistUser(ctx, user.Username)
	if err != nil {
		s.logger.Error("Error during check exist user", zap.String("userName", user.Username), zap.Error(err))
		return err
	}

	if exist {
		return model.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error during generate password hash", zap.String("userName", user.Username), zap.Error(err))
		return err
	}

	err = s.repository.CreateUser(ctx, user.Username, string(hashedPassword))
	if err != nil {
		s.logger.Error("Error during save user", zap.String("userName", user.Username), zap.Error(err))
		return err
	}

	return nil
}

func (s *UserService) FindUser(ctx context.Context, user model.User) (model.User, error) {
	if user.Username == "" || user.Password == "" {
		return model.User{}, model.ErrUserDataIsNotValid
	}

	user, err := s.repository.FindUser(ctx, user.Username)
	if err != nil {
		s.logger.Error("Error during find user", zap.String("userName", user.Username), zap.Error(err))
		return model.User{}, err
	}
	return user, nil
}
