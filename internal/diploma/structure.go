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
	Current  float32 `json:"current"`
	Withdraw float32 `json:"withdraw"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Accural struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
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
	Accural    float32 `json:"accural"`
}
