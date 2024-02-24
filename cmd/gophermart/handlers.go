package main

import (
	"encoding/json"
	"fmt"
	"github.com/gam6itko/go-musthave-diploma/internal"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strings"
)

func decodeLoginPass(body io.ReadCloser) (l *internal.LoginPass, err error) {
	defer body.Close()

	l = new(internal.LoginPass)
	decoder := json.NewDecoder(body)
	err = decoder.Decode(l)
	return
}

// регистрация пользователя;
//
// Возможные коды ответа:
//
//	200 — пользователь успешно зарегистрирован и аутентифицирован;
//	400 — неверный формат запроса;
//	409 — логин уже занят;
//	500 — внутренняя ошибка сервера.
func postUserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusInternalServerError)
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
	userRepo := internal.NewUserRepository(_db)
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
	userId, err := userRepo.InsertNew(r.Context(), *l.Login, string(hashPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//jwt token create
	tokenString, err := _jwtIssuer.IssueFor(userId)
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
func postUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusInternalServerError)
		return
	}

	l, err := decodeLoginPass(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userRepo := internal.NewUserRepository(_db)
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
	tokenString, err := _jwtIssuer.IssueFor(u.Id)
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
func postUserOrders(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		http.Error(w, "Authorization header is empty", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
		return
	}

	_, err := _jwtIssuer.Parse(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

}

// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
func getUserOrders(w http.ResponseWriter, r *http.Request) {

}

// получение текущего баланса счёта баллов лояльности пользователя;
func getUserBalance(w http.ResponseWriter, r *http.Request) {

}

// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
func postUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

}

// получение информации о выводе средств с накопительного счёта пользователем.
func getUserWithdrawals(w http.ResponseWriter, r *http.Request) {

}
