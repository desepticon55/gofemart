package withdrawal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gofemart/internal/model"
	"go.uber.org/zap"
	"net/http"
)

func FindAllWithdrawalsHandler(logger *zap.Logger, service withdrawalService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		withdrawals, err := service.FindAllWithdrawals(request.Context())
		if err != nil {
			if errors.Is(err, model.ErrWithdrawalsWasNotFound) {
				http.Error(writer, "Order number is not filled", http.StatusNoContent)
				return
			}
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(withdrawals)
		if err != nil {
			logger.Error("Error during marshal withdrawals.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err = writer.Write(bytes); err != nil {
			logger.Error("Error write withdrawals.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
