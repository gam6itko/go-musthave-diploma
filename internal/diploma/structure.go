package diploma

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginPass struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

type UserBalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Accrual struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Accrual     float32 `json:"accrual"`
}

type WithdrawalResponse struct {
	Order       string  `json:"order"`
	ProcessedAt string  `json:"processed_at"`
	Sum         float32 `json:"sum"`
}

type OrderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	UploadedAt string  `json:"uploaded_at"`
	Accrual    float32 `json:"accrual"`
}
