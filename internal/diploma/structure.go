package diploma

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

type LoginPass struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

type User struct {
	ID           uint64
	Login        string
	PasswordHash []byte
}

type OrderStatus byte

func OrderStatusFromString(str string) (s OrderStatus, err error) {
	switch strings.ToUpper(str) {
	case "REGISTERED":
		s = StatusRegistered
		return
	case "INVALID":
		s = StatusInvalid
		return
	case "PROCESSING":
		s = StatusProcessing
		return
	case "PROCESSED":
		s = StatusProcessed
		return
	default:
		err = fmt.Errorf("unknown status: %s", str)
	}

	return
}

const (
	// Статуса нет Accural. Это для внутреннего пользования.
	StatusUndefined OrderStatus = iota

	// заказ зарегистрирован, но вознаграждение не рассчитано;
	StatusRegistered

	//заказ не принят к расчёту, и вознаграждение не будет начислено;
	StatusInvalid

	//расчёт начисления в процессе;
	StatusProcessing

	//расчёт начисления окончен;
	StatusProcessed
)

type Order struct {
	ID      uint64
	UserID  uint64
	Status  OrderStatus
	Accural float64
}

type Accural struct {
	OrderNumber string `json:"order"`
	Status      string `json:"status"`
}
