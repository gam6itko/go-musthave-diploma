package diploma

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID              uint64
	Login           string
	PasswordHash    []byte
	BalanceCurrent  float32
	BalanceWithdraw float32
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

func OrderStatusToString(s OrderStatus) (str string, err error) {
	switch s {
	case StatusNew:
		str = "NEW"
		return
	case StatusRegistered:
		str = "REGISTERED"
		return
	case StatusInvalid:
		str = "INVALID"
		return
	case StatusProcessing:
		str = "PROCESSING"
		return
	case StatusProcessed:
		str = "PROCESSED"
		return

	}

	return "", errors.New("unexpected status")
}

const (
	// Статуса нет Accrual. Это для внутреннего пользования.
	StatusNew OrderStatus = iota

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
	ID         uint64
	UserID     uint64
	UploadedAt time.Time
	Status     OrderStatus
	Accrual    float32
}

type Withdrawal struct {
	ID          uint64
	UserID      uint64
	OrderID     uint64
	ProcessedAt time.Time
	Sum         float32
}
