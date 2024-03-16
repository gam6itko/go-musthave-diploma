package controller

import (
	"encoding/json"
	"github.com/gam6itko/go-musthave-diploma/internal/accrual"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/gam6itko/go-musthave-diploma/internal/repository"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type OrderController struct {
	authTrait
	accClient *accrual.Client
	orderRepo *repository.OrderRepository
}

func NewOrderController(accClient *accrual.Client, orderRepo *repository.OrderRepository) *OrderController {
	return &OrderController{
		authTrait{},
		accClient,
		orderRepo,
	}
}

// загрузка пользователем номера заказа для расчёта;
//
// Возможные коды ответа:
// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 401 — пользователь не аутентифицирован;
// 409 — номер заказа уже был загружен другим пользователем;
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func (ths OrderController) PostUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, success := ths.authenticate(w, r)
	if !success {
		return
	}

	// проверка номера
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}
	bNumber, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	orderID, err := strconv.ParseUint(
		strings.Trim(string(bNumber), " \n\t"),
		10,
		64,
	)
	if err != nil {
		http.Error(w, "failed to parse orderID", http.StatusBadRequest)
		return
	}
	if !diploma.LuhnValidate(orderID) {
		http.Error(w, "orderID validation fail", http.StatusUnprocessableEntity)
		return
	}
	orderEntity, err := ths.orderRepo.FindByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "order check fail", http.StatusInternalServerError)
		return
	}
	if orderEntity != nil {
		if orderEntity.UserID == userID {
			http.Error(w, "already processed", http.StatusOK)
		} else {
			http.Error(w, "already processed by another user", http.StatusConflict)
		}
		return
	}

	order := &diploma.Order{
		ID:     orderID,
		UserID: userID,
	}
	if acc, err := ths.accClient.Get(orderID); err != nil {
		log.Printf("failed to get accrual info. for: %d. %s", orderID, err)
	} else {
		if s, sErr := diploma.OrderStatusFromString(acc.Status); sErr != nil {
			log.Printf("failed to get accrual status. %s", sErr)
		} else {
			order.Status = s
			order.Accrual = acc.Accrual
			log.Printf("accrual: %s, %s, %f", acc.OrderNumber, acc.Status, acc.Accrual)
		}
	}

	err = ths.orderRepo.InsertNew(
		r.Context(),
		order,
	)
	if err != nil {
		http.Error(w, "fail to register order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func (ths OrderController) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, success := ths.authenticate(w, r)
	if !success {
		return
	}
	orderList, err := ths.orderRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to get orders", http.StatusInternalServerError)
		return
	}
	if len(orderList) == 0 {
		http.Error(w, "no orders", http.StatusNoContent)
		return
	}

	responseData := make([]*diploma.OrderResponse, len(orderList))
	for i, o := range orderList {
		statusStr, err := diploma.OrderStatusToString(o.Status)
		if err != nil {
			http.Error(w, "no orders", http.StatusInternalServerError)
			return
		}
		responseData[i] = &diploma.OrderResponse{
			Number:     strconv.FormatUint(o.ID, 10),
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
			Status:     statusStr,
			Accrual:    o.Accrual,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Println(err.Error())
	}
}
