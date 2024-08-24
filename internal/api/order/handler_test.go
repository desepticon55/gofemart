package order

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockOrderService struct {
	UploadOrderFunc   func(ctx context.Context, order string) error
	FindAllOrdersFunc func(ctx context.Context) ([]model.Order, error)
}

func (m *mockOrderService) UploadOrder(ctx context.Context, order string) error {
	return m.UploadOrderFunc(ctx, order)
}

func (m *mockOrderService) FindAllOrders(ctx context.Context) ([]model.Order, error) {
	return m.FindAllOrdersFunc(ctx)
}

func TestUploadOrderHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		body           string
		service        orderService
		expectedStatus int
	}{
		{
			name:   "Successful upload",
			method: http.MethodPost,
			body:   "12345",
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return nil
				},
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           "12345",
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Order number is not filled",
			method: http.MethodPost,
			body:   `{}`,
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return model.ErrOrderNumberIsNotFilled
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Order number is not valid",
			method: http.MethodPost,
			body:   "invalid",
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return model.ErrOrderNumberIsNotValid
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "Order number has been uploaded by another user",
			method: http.MethodPost,
			body:   "12345",
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return model.ErrOrderNumberHasUploadedOtherUser
				},
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:   "Order number has been uploaded by the current user",
			method: http.MethodPost,
			body:   "12345",
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return model.ErrOrderNumberHasUploadedCurrentUser
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Internal server error",
			method: http.MethodPost,
			body:   "12345",
			service: &mockOrderService{
				UploadOrderFunc: func(ctx context.Context, order string) error {
					return errors.New("general error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/upload", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := UploadOrderHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestFindAllOrdersHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		service        orderService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful return orders",
			method: http.MethodGet,
			service: &mockOrderService{
				FindAllOrdersFunc: func(ctx context.Context) ([]model.Order, error) {
					return []model.Order{
						{OrderNumber: "12345"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"accrual":0, "number":"12345", "status":"", "uploaded_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Orders not found",
			method: http.MethodGet,
			service: &mockOrderService{
				FindAllOrdersFunc: func(ctx context.Context) ([]model.Order, error) {
					return nil, model.ErrOrdersWasNotFound
				},
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:   "Internal server error",
			method: http.MethodGet,
			service: &mockOrderService{
				FindAllOrdersFunc: func(ctx context.Context) ([]model.Order, error) {
					return nil, errors.New("general error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/orders", nil)
			rec := httptest.NewRecorder()

			handler := FindAllOrdersHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != "" {
				body, err := io.ReadAll(res.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expectedBody, string(body))
			}
		})
	}
}
