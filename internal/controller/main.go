package controller

import (
	"github.com/gam6itko/go-musthave-diploma/internal/jwt"
	"net/http"
	"strings"
)

type authTrait struct {
	jwtIssuer jwt.IIssuerParser
}

// аутентификация
func (ths authTrait) authenticate(w http.ResponseWriter, r *http.Request) (userID uint64, success bool) {
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
	userID, err := ths.jwtIssuer.Parse(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	return userID, true
}
