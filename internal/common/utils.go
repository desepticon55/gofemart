package common

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"hash/fnv"
	"strconv"
	"time"
)

type ContextKey string

const (
	UserNameContextKey ContextKey = "userName"
	Module             int        = 256
)

var JwtKey = []byte("hard_coded_jwt_secret_key")

func CreateJWTToken(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func IsValidOrderNumber(orderNumber string) bool {
	sum := 0
	needDouble := false

	if orderNumber == "" {
		return false
	}

	for i := len(orderNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}
		if needDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		needDouble = !needDouble
	}
	return sum%10 == 0
}

func HashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type TransactionFunc func(tx pgx.Tx) error

func Transactional(ctx context.Context, pool *pgxpool.Pool, fn TransactionFunc) (err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error during open transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)
	return err
}
