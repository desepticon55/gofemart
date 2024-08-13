package orderworker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"github.com/gojek/heimdall/v7/httpclient"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"net/http"
	"strconv"
	"time"
)

type Worker struct {
	from            int
	to              int
	logger          *zap.Logger
	httpClient      *httpclient.Client
	limiter         *rate.Limiter
	orderRepository orderRepository
}

func NewWorker(logger *zap.Logger, repository orderRepository, client *httpclient.Client, from, to int) *Worker {
	limiter := rate.NewLimiter(rate.Limit(10), 1)

	logger.Debug("Make worker", zap.Int("from", from), zap.Int("to", to))
	return &Worker{
		from:            from,
		to:              to,
		httpClient:      client,
		logger:          logger,
		orderRepository: repository,
		limiter:         limiter,
	}
}

func (w *Worker) ProcessOrders(ctx context.Context, accrualAddress string) {
	for {
		orders, err := w.orderRepository.FindOrdersToProcess(ctx, w.from, w.to)
		if err != nil {
			w.logger.Error("Error during fetch orders to process", zap.Error(err))
		}

		for _, order := range orders {
			if err := w.limiter.Wait(ctx); err != nil {
				w.logger.Error("Error during wait rate limiter", zap.Error(err))
			}

			if err := w.processOrder(ctx, accrualAddress, order); err != nil {
				w.logger.Error("Error during process order", zap.Error(err))
			}
		}
	}
}

func (w *Worker) processOrder(ctx context.Context, accrualAddress string, order common.Order) error {
	url := fmt.Sprintf("%s/api/orders/%s", accrualAddress, order.OrderNumber)
	w.logger.Debug("Accrual address prepared", zap.String("address", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error during create request: %w", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error during send event: %w", err)
	}
	defer resp.Body.Close()

	var accrual struct {
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual"`
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			retryDelay, err := strconv.Atoi(retryAfter)
			if err != nil {
				retryDelay = 5
			}
			w.logger.Debug(fmt.Sprintf("Received 429, retrying after %d seconds", retryDelay))
			time.Sleep(time.Duration(retryDelay) * time.Second)
			return w.processOrder(ctx, accrualAddress, order)
		}
	}

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&accrual)
		if err != nil {
			return fmt.Errorf("error during decode response: %w", err)
		}
		w.logger.Debug(
			"Received accrual response",
			zap.String("order", accrual.Order),
			zap.String("status", accrual.Status),
			zap.Float64("accrual", accrual.Accrual))

		err := w.orderRepository.ChangeOrderStatus(ctx, order, accrual.Status, accrual.Accrual)
		if err != nil {
			return fmt.Errorf("error during chage order: %w", err)
		}
	}

	return nil
}
