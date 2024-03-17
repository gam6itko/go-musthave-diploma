package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

type IIssuer interface {
	Issue(userID uint64) (tokenString string, err error)
}

type IParser interface {
	Parse(tokenString string) (uint64, error)
}

type IIssuerParser interface {
	IIssuer
	IParser
}

type Issuer struct {
	key []byte
}

func NewIssuer(key []byte) *Issuer {
	return &Issuer{
		key,
	}
}

func (ths Issuer) Issue(userID uint64) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: &jwt.NumericDate{
					Time: time.Now().UTC(),
				},
				ExpiresAt: &jwt.NumericDate{
					Time: time.Now().Add(24 * time.Hour).UTC(),
				},
			},
			UserID: userID,
		},
	)

	// создаём строку токена
	tokenString, err = token.SignedString(ths.key)
	return
}

func (ths Issuer) Parse(tokenString string) (uint64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return ths.key, nil
		},
	)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}
