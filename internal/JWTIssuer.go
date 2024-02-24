package internal

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTIssuer struct {
	key []byte
}

func NewJWTIssuer(key []byte) *JWTIssuer {
	return &JWTIssuer{
		key,
	}
}

func (ths JWTIssuer) IssueFor(userId int64) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: &jwt.NumericDate{
					Time: time.Now().UTC(),
				},
			},
			UserID: userId,
		},
	)

	// создаём строку токена
	tokenString, err = token.SignedString(ths.key)
	return
}
