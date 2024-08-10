package common

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	NewOrderStatus        = "NEW"
	ProcessingOrderStatus = "PROCESSING"
	InvalidOrderStatus    = "INVALID"
	ProcessedOrderStatus  = "PROCESSED"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Username string
	Balance  float64
	Version  int64
}

type BalanceStats struct {
	Username  string  `json:"-"`
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	ID          string
	Username    string
	OrderNumber string
	Sum         float64
	CreateDate  time.Time
}

func (e *Withdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		OrderNumber string  `json:"order"`
		Sum         float64 `json:"sum"`
		CreateDate  string  `json:"processed_at"`
	}{
		OrderNumber: e.OrderNumber,
		CreateDate:  e.CreateDate.Format(time.RFC3339),
		Sum:         e.Sum,
	})
}

type Order struct {
	OrderNumber    string
	CreateDate     time.Time
	LastModifyDate time.Time
	Status         string
	Username       string
	Accrual        float64
	KeyHash        int64
	KeyHashModule  int64
	Version        int64
}

func (e *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		OrderNumber string  `json:"number"`
		CreateDate  string  `json:"uploaded_at"`
		Status      string  `json:"status"`
		Accrual     float64 `json:"accrual"`
	}{
		OrderNumber: e.OrderNumber,
		CreateDate:  e.CreateDate.Format(time.RFC3339),
		Status:      e.Status,
		Accrual:     e.Accrual,
	})
}
