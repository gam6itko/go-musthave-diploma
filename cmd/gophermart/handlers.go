package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma/repository"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func decodeLoginPass(body io.ReadCloser) (l *diploma.LoginPass, err error) {
	defer body.Close()

	l = new(diploma.LoginPass)
	decoder := json.NewDecoder(body)
	err = decoder.Decode(l)
	return
}

// аутентификация
func authenticate(w http.ResponseWriter, r *http.Request) (userID uint64, success bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
		return
	}
	userID, err := _jwtIssuer.Parse(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	return userID, true
}

// регистрация пользователя;
//
// Возможные коды ответа:
//
//	200 — пользователь успешно зарегистрирован и аутентифицирован;
//	400 — неверный формат запроса;
//	409 — логин уже занят;
//	500 — внутренняя ошибка сервера.
func postUserRegister(w http.ResponseWriter, r *http.Request, db *sql.DB, jwtIssuer *diploma.JWTIssuer) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	l, err := decodeLoginPass(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if l.Login == nil || l.Password == nil || len(*l.Login) == 0 || len(*l.Password) == 0 {
		http.Error(w, "missing required credentials", http.StatusBadRequest)
		return
	}
	if *l.Password == "123" {
		http.Error(w, "password is to weak. try 'qwerty' ^_^", http.StatusBadRequest)
		return
	}

	// check user exists
	userRepo := repository.NewUserRepository(db)
	u, err := userRepo.FindByLogin(r.Context(), *l.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u != nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	// user save
	hashPass, err := bcrypt.GenerateFromPassword([]byte(*l.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID, err := userRepo.InsertNew(r.Context(), *l.Login, string(hashPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//jwt token create
	tokenString, err := jwtIssuer.Issue(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(200)
}

// аутентификация пользователя;
//
// Возможные коды ответа:
// 200 — пользователь успешно аутентифицирован;
// 400 — неверный формат запроса;
// 401 — неверная пара логин/пароль;
// 500 — внутренняя ошибка сервера.
func postUserLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, jwtIssuer *diploma.JWTIssuer) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	l, err := decodeLoginPass(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userRepo := repository.NewUserRepository(db)
	u, err := userRepo.FindByLogin(r.Context(), *l.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if u == nil {
		// wrong user or password
		http.Error(w, "user does not exists ^_^", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(*l.Password)); err != nil {
		// wrong user or password
		http.Error(w, "wrong user password ^_^", http.StatusUnauthorized)
		return
	}

	//jwt issue
	tokenString, err := jwtIssuer.Issue(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(200)
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
func postUserOrders(w http.ResponseWriter, r *http.Request, db *sql.DB, accClient *diploma.AccuralClient) {
	userID, success := authenticate(w, r)
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
	repo := repository.NewOrderRepository(db)
	orderEntity, err := repo.FindByID(r.Context(), orderID)
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
	if acc, err := accClient.Get(orderID); err != nil {
		log.Printf("failed to get accural info. %s", err)
	} else {
		if s, sErr := diploma.OrderStatusFromString(acc.Status); sErr != nil {
			log.Printf("failed to get accural status. %s", sErr)
		} else {
			order.Status = s
			order.Accural = acc.Accrual
			log.Printf("accural: %s, %s, %f", acc.OrderNumber, acc.Status, acc.Accrual)
		}
	}

	err = repo.InsertNew(
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
func getUserOrders(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, success := authenticate(w, r)
	if !success {
		return
	}
	repo := repository.NewOrderRepository(db)
	orderList, err := repo.FindByUserID(r.Context(), userID)
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
			http.Error(w, "no orders", http.StatusNoContent)
			return
		}
		responseData[i] = &diploma.OrderResponse{
			Number:     strconv.FormatUint(o.ID, 10),
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
			Status:     statusStr,
			Accural:    o.Accural,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		log.Println(err.Error())
	}
}

// получение текущего баланса счёта баллов лояльности пользователя;
func getUserBalance(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, success := authenticate(w, r)
	if !success {
		return
	}
	repo := repository.NewUserRepository(db)
	u, err := repo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to retrieve user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&diploma.UserBalanceResponse{
		Current:  u.BalanceCurrent,
		Withdraw: u.BalanceWithdraw,
	})
	if err != nil {
		log.Printf("encode error: %s", err)
	}
}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
func postUserBalanceWithdraw(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, success := authenticate(w, r)
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

	repo := repository.NewUserRepository(db)
	u, err := repo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to retrieve user data", http.StatusInternalServerError)
		return
	}
	if u.BalanceCurrent < reqData.Sum {
		http.Error(w, "not enough balance", http.StatusPaymentRequired)
		return
	}

	if err := repo.Withdraw(r.Context(), userID, orderID, reqData.Sum); err != nil {
		http.Error(w, "fail to withdraw", http.StatusInternalServerError)
		return
	}
}

// получение информации о выводе средств с накопительного счёта пользователем.
func getUserWithdrawals(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, success := authenticate(w, r)
	if !success {
		return
	}

	repo := repository.NewWithdrawalRepository(db)
	wList, err := repo.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to get withdrawal list", http.StatusInternalServerError)
		return
	}

	responseData := make([]*diploma.WithdrawalResponse, len(wList))
	for i, w := range wList {
		responseData[i] = &diploma.WithdrawalResponse{
			Order:       strconv.FormatUint(w.OrderID, 10),
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
			Sum:         w.Sum,
		}
	}
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(responseData); err != nil {
		log.Printf("encode error: %s", err)
	}
}
