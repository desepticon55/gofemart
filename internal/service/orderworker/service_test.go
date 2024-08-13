package orderworker

import (
	"context"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) FindOrdersToProcess(ctx context.Context, from, to int) ([]common.Order, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]common.Order), args.Error(1)
}

func (m *MockOrderRepository) ChangeOrderStatus(ctx context.Context, order common.Order, status string, accrual float64) error {
	args := m.Called(ctx, order, status, accrual)
	return args.Error(0)
}

func TestWorker_ProcessOrders(t *testing.T) {
	logger := zaptest.NewLogger(t)
	limiter := rate.NewLimiter(10, 1)

	t.Run("should successfully processed order status if response status is 200", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		client := httpclient.NewClient(httpclient.WithHTTPTimeout(10 * time.Millisecond))
		mockRepo := new(MockOrderRepository)
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "order": "12345", "status": "COMPLETED", "accrual": 100 }`))
		}
		server := httptest.NewServer(http.HandlerFunc(dummyHandler))
		defer server.Close()

		var worker = &Worker{
			from:            0,
			to:              10,
			logger:          logger,
			httpClient:      client,
			limiter:         limiter,
			orderRepository: mockRepo,
		}
		order := common.Order{OrderNumber: "12345"}
		mockRepo.On("FindOrdersToProcess", ctx, 0, 10).Return([]common.Order{order}, nil).Once().On("FindOrdersToProcess", ctx, 0, 10).Return([]common.Order{}, nil)
		mockRepo.On("ChangeOrderStatus", ctx, order, "COMPLETED", 100.0).Return(nil)

		go worker.ProcessOrders(ctx, server.URL)

		<-ctx.Done()

		mockRepo.AssertCalled(t, "ChangeOrderStatus", ctx, order, "COMPLETED", 100.0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should successfully process order status if response status is 200 after 429", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var callCount int
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			callCount++
			if callCount == 1 {
				w.Header().Set("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "order": "12345", "status": "COMPLETED", "accrual": 100 }`))
		}
		server := httptest.NewServer(http.HandlerFunc(dummyHandler))
		defer server.Close()

		client := httpclient.NewClient(httpclient.WithHTTPTimeout(10 * time.Millisecond))
		mockRepo := new(MockOrderRepository)

		var worker = &Worker{
			from:            0,
			to:              10,
			logger:          logger,
			httpClient:      client,
			limiter:         limiter,
			orderRepository: mockRepo,
		}
		order := common.Order{OrderNumber: "12345"}
		mockRepo.On("FindOrdersToProcess", ctx, 0, 10).Return([]common.Order{order}, nil).Once().On("FindOrdersToProcess", ctx, 0, 10).Return([]common.Order{}, nil)
		mockRepo.On("ChangeOrderStatus", ctx, order, "COMPLETED", 100.0).Return(nil)

		go worker.ProcessOrders(ctx, server.URL)

		<-ctx.Done()

		mockRepo.AssertCalled(t, "ChangeOrderStatus", ctx, order, "COMPLETED", 100.0)
		mockRepo.AssertExpectations(t)
	})
}
