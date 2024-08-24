package user

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) ExistUser(ctx context.Context, userName string) (bool, error) {
	args := m.Called(ctx, userName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, userName, password string) error {
	args := m.Called(ctx, userName, password)
	return args.Error(0)
}

func (m *MockUserRepository) FindUser(ctx context.Context, userName string) (model.User, error) {
	args := m.Called(ctx, userName)
	return args.Get(0).(model.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")
	logger := zaptest.NewLogger(t)

	t.Run("should successfully create user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		user := model.User{Username: "newUser", Password: "password"}
		mockRepo.On("ExistUser", ctx, "newUser").Return(false, nil)
		mockRepo.On("CreateUser", ctx, "newUser", mock.Anything).Return(nil)

		err := service.CreateUser(ctx, user)
		assert.NoError(t, err)

		mockRepo.AssertCalled(t, "ExistUser", ctx, "newUser")
		mockRepo.AssertCalled(t, "CreateUser", ctx, "newUser", mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if user exist", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		user := model.User{Username: "newUser", Password: "password"}
		mockRepo.On("ExistUser", ctx, "newUser").Return(true, nil)

		err := service.CreateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, err, model.ErrUserAlreadyExists)

		mockRepo.AssertCalled(t, "ExistUser", ctx, "newUser")
		mockRepo.AssertNotCalled(t, "CreateUser", ctx, "newUser", mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if ExistUser(..) return error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		expectedError := errors.New("database error")
		user := model.User{Username: "newUser", Password: "password"}
		mockRepo.On("ExistUser", ctx, "newUser").Return(false, expectedError)

		err := service.CreateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, err, expectedError)

		mockRepo.AssertCalled(t, "ExistUser", ctx, "newUser")
		mockRepo.AssertNotCalled(t, "CreateUser", ctx, "newUser", mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if CreateUser(..) return error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		expectedError := errors.New("database error")
		user := model.User{Username: "newUser", Password: "password"}
		mockRepo.On("ExistUser", ctx, "newUser").Return(false, nil)
		mockRepo.On("CreateUser", ctx, "newUser", mock.Anything).Return(expectedError)

		err := service.CreateUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, err, expectedError)

		mockRepo.AssertCalled(t, "ExistUser", ctx, "newUser")
		mockRepo.AssertCalled(t, "CreateUser", ctx, "newUser", mock.Anything)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_FindUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")
	logger := zaptest.NewLogger(t)

	t.Run("should return found user", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		user := model.User{Username: "newUser", Password: "password"}
		mockRepo.On("FindUser", ctx, "newUser").Return(user, nil)

		foundUser, err := service.FindUser(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, user, foundUser)

		mockRepo.AssertCalled(t, "FindUser", ctx, "newUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if FindUser(..) return error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		user := model.User{Username: "newUser", Password: "password"}
		expectedError := errors.New("database error")
		mockRepo.On("FindUser", ctx, "newUser").Return(model.User{}, expectedError)

		_, err := service.FindUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockRepo.AssertCalled(t, "FindUser", ctx, "newUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if login or password is empty", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := &UserService{
			repository: mockRepo,
			logger:     logger,
		}

		user := model.User{Username: "", Password: ""}

		_, err := service.FindUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, model.ErrUserDataIsNotValid, err)

		mockRepo.AssertNotCalled(t, "FindUser", ctx, "newUser")
		mockRepo.AssertExpectations(t)
	})
}
