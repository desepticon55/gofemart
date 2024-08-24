package withdrawal

import (
	"context"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockWithdrawalService struct {
	FindAllWithdrawalsFunc func(ctx context.Context) ([]model.Withdrawal, error)
}

func (m *mockWithdrawalService) FindAllWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	return m.FindAllWithdrawalsFunc(ctx)
}

func TestFindAllWithdrawalsHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		service        withdrawalService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful return withdrawals",
			method: http.MethodGet,
			service: &mockWithdrawalService{
				FindAllWithdrawalsFunc: func(ctx context.Context) ([]model.Withdrawal, error) {
					return []model.Withdrawal{
						{ID: "1", Sum: 100.0, Username: "testUser", OrderNumber: "12345"},
						{ID: "2", Sum: 50.0, Username: "testUser", OrderNumber: "67432"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"order":"12345", "processed_at":"0001-01-01T00:00:00Z", "sum":100}, {"order":"67432", "processed_at":"0001-01-01T00:00:00Z", "sum":50}]`,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			service:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Withdrawals not found",
			method: http.MethodGet,
			service: &mockWithdrawalService{
				FindAllWithdrawalsFunc: func(ctx context.Context) ([]model.Withdrawal, error) {
					return nil, model.ErrWithdrawalsWasNotFound
				},
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:   "Internal server error",
			method: http.MethodGet,
			service: &mockWithdrawalService{
				FindAllWithdrawalsFunc: func(ctx context.Context) ([]model.Withdrawal, error) {
					return nil, errors.New("general error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/withdrawals", nil)
			rec := httptest.NewRecorder()

			handler := FindAllWithdrawalsHandler(logger, tt.service)
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
