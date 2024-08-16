package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gofemart/internal/model"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func UploadOrderHandler(logger *zap.Logger, service orderService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			logger.Error("Error during read body", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()

		logger.Debug(fmt.Sprintf("Received body = %s", body))

		err = service.UploadOrder(request.Context(), string(body))
		if err != nil {
			if errors.Is(err, model.ErrOrderNumberIsNotFilled) {
				http.Error(writer, "Order number is not filled", http.StatusBadRequest)
				return
			}

			if errors.Is(err, model.ErrOrderNumberIsNotValid) {
				http.Error(writer, "Order number is not valid", http.StatusUnprocessableEntity)
				return
			}

			if errors.Is(err, model.ErrOrderNumberHasUploadedOtherUser) {
				http.Error(writer, "Order number has uploaded from other user", http.StatusConflict)
				return
			}

			if errors.Is(err, model.ErrOrderNumberHasUploadedCurrentUser) {
				writer.WriteHeader(http.StatusOK)
				return
			}

			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusAccepted)
	}
}

func FindAllOrdersHandler(logger *zap.Logger, service orderService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		orders, err := service.FindAllOrders(request.Context())
		if err != nil {
			if errors.Is(err, model.ErrOrdersWasNotFound) {
				http.Error(writer, "Orders was not found", http.StatusNoContent)
				return
			}
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(orders)
		if err != nil {
			logger.Error("Error during marshal orders.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if _, err = writer.Write(bytes); err != nil {
			logger.Error("Error write orders.", zap.Error(err))
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
