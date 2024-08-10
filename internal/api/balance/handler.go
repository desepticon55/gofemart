package balance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"go.uber.org/zap"
	"net/http"
)

func FindUserBalanceHandler(logger *zap.Logger, service balanceService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		balance, err := service.FindBalanceStats(request.Context())
		if err != nil {
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(balance)
		if err != nil {
			logger.Error("Error during marshal balance.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err = writer.Write(bytes); err != nil {
			logger.Error("Error write balance.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func WithdrawBalanceHandler(logger *zap.Logger, service balanceService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		var req struct {
			OrderNumber string  `json:"order"`
			Sum         float64 `json:"sum"`
		}

		err := json.NewDecoder(request.Body).Decode(&req)
		if err != nil {
			logger.Error("Invalid request payload", zap.Error(err))
			http.Error(writer, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = service.Withdraw(request.Context(), req.OrderNumber, req.Sum)
		if err != nil {
			if errors.Is(err, common.ErrOrderNumberOrSumIsNotFilled) {
				http.Error(writer, "Order number or sum is not filled", http.StatusBadRequest)
				return
			}

			if errors.Is(err, common.ErrOrderNumberIsNotValid) {
				http.Error(writer, "Order number is not valid", http.StatusUnprocessableEntity)
				return
			}

			if errors.Is(err, common.ErrUserBalanceLessThanSumToWithdraw) {
				http.Error(writer, "Balance less than sum to withdrawal", http.StatusPaymentRequired)
				return
			}

			if errors.Is(err, common.ErrUserBalanceHasChanged) {
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
