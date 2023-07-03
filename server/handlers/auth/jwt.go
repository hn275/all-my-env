package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

func Sign(data jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secret))
}

func Decode(d string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(d, &JwtToken{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtToken)
	if !ok || !token.Valid {
		return nil, ErrInvalidType
	}

	return claims, nil
}
