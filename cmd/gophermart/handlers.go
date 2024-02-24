package main

import (
	"encoding/json"
	"fmt"
	"github.com/gam6itko/go-musthave-diploma/internal"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// регистрация пользователя;
// Статусы ответа:
//
//	200 — пользователь успешно зарегистрирован и аутентифицирован;
//	400 — неверный формат запроса;
//	409 — логин уже занят;
//	500 — внутренняя ошибка сервера.
func postUserRegister(w http.ResponseWriter, r *http.Request) {
	l := new(internal.LoginPass)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	exists, err := userRepo.IsExists(r.Context(), *l.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	// user save
	hashPass, err := bcrypt.GenerateFromPassword([]byte(*l.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := userRepo.InsertNew(r.Context(), *l.Login, string(hashPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//jwt token create
	jwtKey, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		jwtKey = "qwerty"
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		internal.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: &jwt.NumericDate{
					Time: time.Now().UTC(),
				},
			},
			UserID: id,
		},
	)

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(200)
}

// аутентификация пользователя;
func postUserLogin(w http.ResponseWriter, r *http.Request) {

}

// загрузка пользователем номера заказа для расчёта;
func postUserOrders(w http.ResponseWriter, r *http.Request) {

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
