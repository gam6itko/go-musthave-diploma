package internal

import "github.com/golang-jwt/jwt/v5"

type LoginPass struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}
