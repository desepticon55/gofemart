package middleware

import (
	"compress/gzip"
	"context"
	"github.com/desepticon55/gofemart/internal/api/auth"
	"github.com/desepticon55/gofemart/internal/model"
	"github.com/desepticon55/gofemart/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

func CheckAuthMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			const bearerPrefix = "Bearer "
			authHeader := request.Header.Get("Authorization")
			if authHeader == "" {
				logger.Error("Authorization header is missing", zap.String("Authorization", authHeader))
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				logger.Error("Invalid Authorization header format", zap.String("Authorization", authHeader))
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenStr := authHeader[len(bearerPrefix):]
			claims := &model.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return auth.JwtKey, nil
			})

			if err != nil {
				logger.Error("Error during parse token", zap.String("Authorization", authHeader), zap.Error(err))
				http.Error(writer, "Expired token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				logger.Error("Invalid token", zap.String("Authorization", authHeader))
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(request.Context(), service.UserNameContextKey, claims.Username)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func DecompressingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.Contains(request.Header.Get("Content-Encoding"), "gzip") {
				reader, err := gzip.NewReader(request.Body)
				if err != nil {
					http.Error(writer, "Error during create gzip reader", http.StatusInternalServerError)
					return
				}
				defer reader.Close()
				request.Body = reader
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func CompressingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			acceptEncoding := request.Header.Get("Accept-Encoding")
			if strings.Contains(acceptEncoding, "gzip") {
				writer.Header().Set("Content-Encoding", "gzip")
				gzipWriter := gzip.NewWriter(writer)
				defer gzipWriter.Close()
				gzipResponseWriter := &gzipResponseWriter{gzipWriter, writer}
				next.ServeHTTP(gzipResponseWriter, request)
			} else {
				next.ServeHTTP(writer, request)
			}
		})
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
