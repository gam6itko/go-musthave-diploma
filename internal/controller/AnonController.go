package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/gam6itko/go-musthave-diploma/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

type AnonController struct {
	jwtIssuer *diploma.JWTIssuer
	userRepo  *repository.UserRepository
}

func (ths AnonController) decodeLoginPass(body io.ReadCloser) (l *diploma.LoginPass, err error) {
	defer body.Close()

	l = new(diploma.LoginPass)
	decoder := json.NewDecoder(body)
	err = decoder.Decode(l)
	return
}

func NewAnonController(jwtIssuer *diploma.JWTIssuer, userRepo *repository.UserRepository) *AnonController {
	return &AnonController{
		jwtIssuer,
		userRepo,
	}
}

// регистрация пользователя;
//
// Возможные коды ответа:
//
//	200 — пользователь успешно зарегистрирован и аутентифицирован;
//	400 — неверный формат запроса;
//	409 — логин уже занят;
//	500 — внутренняя ошибка сервера.
func (ths AnonController) PostUserRegister(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	l, err := ths.decodeLoginPass(r.Body)
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
	u, err := ths.userRepo.FindByLogin(r.Context(), *l.Login)
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
	userID, err := ths.userRepo.InsertNew(r.Context(), *l.Login, string(hashPass))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//jwt token create
	tokenString, err := ths.jwtIssuer.Issue(userID)
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
func (ths AnonController) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	l, err := ths.decodeLoginPass(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := ths.userRepo.FindByLogin(r.Context(), *l.Login)
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
	tokenString, err := ths.jwtIssuer.Issue(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(200)
}
