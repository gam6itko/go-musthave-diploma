package controller

import (
	"encoding/json"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"github.com/gam6itko/go-musthave-diploma/internal/repository"
	"log"
	"net/http"
)

type UserController struct {
	authTrait
	userRepo *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{
		authTrait{},
		userRepo,
	}
}

// получение текущего баланса счёта баллов лояльности пользователя;
func (ths UserController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, success := ths.authenticate(w, r)
	if !success {
		return
	}
	u, err := ths.userRepo.FindByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "fail to retrieve user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(
		&diploma.UserBalanceResponse{
			Current:   u.BalanceCurrent,
			Withdrawn: u.BalanceWithdraw,
		},
	)
	if err != nil {
		log.Printf("encode error: %s", err)
	}
}
