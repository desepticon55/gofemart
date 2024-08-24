package auth

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockUserService struct {
	FindUserFunc   func(ctx context.Context, user model.User) (model.User, error)
	CreateUserFunc func(ctx context.Context, user model.User) error
}

func (m *mockUserService) FindUser(ctx context.Context, user model.User) (model.User, error) {
	return m.FindUserFunc(ctx, user)
}

func (m *mockUserService) CreateUser(ctx context.Context, user model.User) error {
	return m.CreateUserFunc(ctx, user)
}

func TestLoginHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	mockUser := model.User{Username: "testUser", Password: string(passwordHash)}

	tests := []struct {
		name           string
		method         string
		body           string
		service        userService
		expectedStatus int
	}{
		{
			name:   "Successful login",
			method: http.MethodPost,
			body:   `{"login":"testUser", "password":"password"}`,
			service: &mockUserService{
				FindUserFunc: func(ctx context.Context, user model.User) (model.User, error) {
					return mockUser, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           "",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty username or password",
			method: http.MethodPost,
			body:   `{"login":"", "password":""}`,
			service: &mockUserService{
				FindUserFunc: func(ctx context.Context, user model.User) (model.User, error) {
					return model.User{}, model.ErrUserDataIsNotValid
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "User not found",
			method: http.MethodPost,
			body:   `{"login":"nonexistent", "password":"password"}`,
			service: &mockUserService{
				FindUserFunc: func(ctx context.Context, user model.User) (model.User, error) {
					return model.User{}, errors.New("user not found")
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Incorrect password",
			method: http.MethodPost,
			body:   `{"login":"testUser", "password":"wrongpassword"}`,
			service: &mockUserService{
				FindUserFunc: func(ctx context.Context, user model.User) (model.User, error) {
					return mockUser, nil
				},
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := LoginHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if res.Status == "OK" {
				auth := res.Header.Get("Authorization")
				assert.NotEmpty(t, auth)
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		body           string
		service        userService
		expectedStatus int
	}{
		{
			name:   "Successful registration",
			method: http.MethodPost,
			body:   `{"login":"testUser", "password":"password"}`,
			service: &mockUserService{
				CreateUserFunc: func(ctx context.Context, user model.User) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           "",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Invalid request payload",
			method: http.MethodPost,
			body:   `{"login":"", "password":""}`,
			service: &mockUserService{
				CreateUserFunc: func(ctx context.Context, user model.User) error {
					return model.ErrUserDataIsNotValid
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "User already exists",
			method: http.MethodPost,
			body:   `{"login":"existingUser", "password":"password"}`,
			service: &mockUserService{
				CreateUserFunc: func(ctx context.Context, user model.User) error {
					return model.ErrUserAlreadyExists
				},
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := RegisterHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if res.Status == "OK" {
				auth := res.Header.Get("Authorization")
				assert.NotEmpty(t, auth)
			}
		})
	}
}
