package controller

import (
	"encoding/json"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	repository2 "github.com/gam6itko/go-musthave-diploma/internal/repository"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WithdrawalController struct {
	authTrait
	wRepo    *repository2.WithdrawalRepository
	userRepo *repository2.UserRepository
}

func NewWithdrawalController(wRepo *repository2.WithdrawalRepository, userRepo *repository2.UserRepository) *WithdrawalController {
	return &WithdrawalController{
		authTrait{},
		wRepo,
		userRepo,
	}
}

// Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
func (ths WithdrawalController) PostUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, success := ths.authenticate(w, r)
	if !success {
		return
	}

	reqData := new(diploma.WithdrawRequest)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(reqData); err != nil {
		http.Error(w, "fail to decode request body", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.ParseUint(reqData.Order, 10, 64)
	if err != nil || !diploma.LuhnValidate(orderID) {
		http.Error(w, "orderID validation fail", http.StatusUnprocessableEntity)
		return
	}

	u, err := ths.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to retrieve user data", http.StatusInternalServerError)
		return
	}
	if u.BalanceCurrent < reqData.Sum {
		http.Error(w, "not enough balance", http.StatusPaymentRequired)
		return
	}

	if err := ths.userRepo.Withdraw(r.Context(), userID, orderID, reqData.Sum); err != nil {
		http.Error(w, "fail to withdraw", http.StatusInternalServerError)
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем.
func (ths WithdrawalController) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, success := ths.authenticate(w, r)
	if !success {
		return
	}

	wList, err := ths.wRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to get withdrawal list", http.StatusInternalServerError)
		return
	}

	responseData := make([]*diploma.WithdrawalResponse, len(wList))
	for i, wt := range wList {
		responseData[i] = &diploma.WithdrawalResponse{
			Order:       strconv.FormatUint(wt.OrderID, 10),
			ProcessedAt: wt.ProcessedAt.Format(time.RFC3339),
			Sum:         wt.Sum,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(responseData); err != nil {
		log.Printf("encode error: %s", err)
	}
}
