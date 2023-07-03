package jsonwebtoken

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret string
)

func init() {
	secret = os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}
}

func Sign(data jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secret))
}
