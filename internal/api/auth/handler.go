package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/desepticon55/gofemart/internal/common"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func LoginHandler(logger *zap.Logger, service userService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		var user common.User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			http.Error(writer, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			http.Error(writer, "Invalid request payload", http.StatusBadRequest)
			return
		}

		foundUser, err := service.FindUser(request.Context(), user.Username)
		if err != nil {
			http.Error(writer, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if err != nil {
			http.Error(writer, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := common.CreateJWTToken(user.Username)
		if err != nil {
			logger.Error("Error during create token", zap.String("username", user.Username), zap.Error(err))
			http.Error(writer, "Could not create token", http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully create token", zap.String("username", user.Username), zap.String("token", token))

		writer.Header().Set("Authorization", "Bearer "+token)
		writer.WriteHeader(http.StatusOK)
	}
}

func RegisterHandler(logger *zap.Logger, service userService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, fmt.Sprintf("Method '%s' is not allowed", request.Method), http.StatusBadRequest)
			return
		}

		var user common.User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			http.Error(writer, "Invalid request payload", http.StatusInternalServerError)
			return
		}

		err = service.CreateUser(request.Context(), user)
		if err != nil {
			if errors.Is(err, common.ErrUserDataIsNotValid) {
				http.Error(writer, "Invalid request payload", http.StatusBadRequest)
				return
			}

			if errors.Is(err, common.ErrUserAlreadyExists) {
				http.Error(writer, fmt.Sprintf("User with login = %s already exist", user.Username), http.StatusConflict)
				return
			}

			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully save user", zap.String("username", user.Username))

		token, err := common.CreateJWTToken(user.Username)
		if err != nil {
			logger.Error("Error during create token", zap.String("username", user.Username), zap.Error(err))
			http.Error(writer, "Could not create token", http.StatusInternalServerError)
			return
		}
		logger.Debug("Successfully create token", zap.String("username", user.Username), zap.String("token", token))

		writer.Header().Set("Authorization", "Bearer "+token)
		writer.WriteHeader(http.StatusOK)
	}
}
