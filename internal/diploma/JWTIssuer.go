package diploma

import (
	"errors"
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

func (ths JWTIssuer) Issue(userID uint64) (tokenString string, err error) {
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

func (ths JWTIssuer) Parse(tokenString string) (uint64, error) {
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
