package jsonwebtoken

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/lib"
)

var (
	secret  string
	Decoder JsonWebToken
)

type JsonWebToken interface {
	Decode(string) (*JwtToken, error)
}

type JwtDecoder struct{}

func init() {
	secret = lib.Getenv("JWT_SECRET")
	Decoder = &JwtDecoder{}
}

type GithubUser struct {
	Token     string `json:"token,omitempty"`
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type JwtToken struct {
	GithubUser `json:",inline"`
	jwt.RegisteredClaims
}

func Sign(data jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secret))
}

func NewDecoder() JsonWebToken {
	return Decoder
}

func (d *JwtDecoder) Decode(t string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(t, &JwtToken{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("invalid signing algo")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtToken)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt type")
	}

	return claims, nil
}

func GetUser(r *http.Request) (*GithubUser, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return nil, errors.New("auth token not found")
	}

	t := strings.Split(h, " ")

	if len(t) != 2 {
		return nil, errors.New("invalid auth token")
	}

	if strings.ToLower(t[0]) != "bearer" {
		return nil, errors.New("invalid auth token")
	}

	decoded, err := Decoder.Decode(t[1])
	if err != nil {
		return nil, err
	}

	return &decoded.GithubUser, nil
}
