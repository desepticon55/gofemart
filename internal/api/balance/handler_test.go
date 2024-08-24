package balance

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockBalanceService struct {
	FindBalanceStatsFunc func(ctx context.Context) (model.BalanceStats, error)
	WithdrawFunc         func(ctx context.Context, orderNumber string, sum float64) error
}

func (m *mockBalanceService) FindBalanceStats(ctx context.Context) (model.BalanceStats, error) {
	return m.FindBalanceStatsFunc(ctx)
}

func (m *mockBalanceService) Withdraw(ctx context.Context, orderNumber string, sum float64) error {
	return m.WithdrawFunc(ctx, orderNumber, sum)
}

func TestFindUserBalanceHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		service        balanceService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Successful balance retrieval",
			method: http.MethodGet,
			service: &mockBalanceService{
				FindBalanceStatsFunc: func(ctx context.Context) (model.BalanceStats, error) {
					return model.BalanceStats{
						Username:  "testUser",
						Balance:   1000.0,
						Withdrawn: 200.0,
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"current":1000,"withdrawn":200}`,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			service:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:   "Internal server error",
			method: http.MethodGet,
			service: &mockBalanceService{
				FindBalanceStatsFunc: func(ctx context.Context) (model.BalanceStats, error) {
					return model.BalanceStats{}, errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), service.UserNameContextKey, "testUser")
			req := httptest.NewRequest(tt.method, "/balance", nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			handler := FindUserBalanceHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != "" {
				var body model.BalanceStats
				err := json.NewDecoder(res.Body).Decode(&body)
				assert.NoError(t, err)

				var expectedBody model.BalanceStats
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedBody)
				assert.NoError(t, err)

				assert.Equal(t, expectedBody, body)
			}
		})
	}
}

func TestWithdrawBalanceHandler(t *testing.T) {
	logger := zaptest.NewLogger(t)
	defer logger.Sync()

	tests := []struct {
		name           string
		method         string
		body           string
		service        balanceService
		expectedStatus int
	}{
		{
			name:   "Successful withdrawal",
			method: http.MethodPost,
			body:   `{"order":"12345","sum":200}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
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
			body:   `{"order":"12345"}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
					return model.ErrOrderNumberOrSumIsNotFilled
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Order number not valid",
			method: http.MethodPost,
			body:   `{"order":"invalid","sum":200}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
					return model.ErrOrderNumberIsNotValid
				},
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "Insufficient balance",
			method: http.MethodPost,
			body:   `{"order":"12345","sum":2000}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
					return model.ErrUserBalanceLessThanSumToWithdraw
				},
			},
			expectedStatus: http.StatusPaymentRequired,
		},
		{
			name:   "Balance has changed",
			method: http.MethodPost,
			body:   `{"order":"12345","sum":200}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
					return model.ErrUserBalanceHasChanged
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "General error",
			method: http.MethodPost,
			body:   `{"order":"12345","sum":200}`,
			service: &mockBalanceService{
				WithdrawFunc: func(ctx context.Context, orderNumber string, sum float64) error {
					return errors.New("general error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/withdraw", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler := WithdrawBalanceHandler(logger, tt.service)
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
