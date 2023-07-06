package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/lib"
)

type AuthCx interface {
	Do(r *http.Request) (*http.Response, error)
}

type Token struct {
	Code string
}

type GithubAuthToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
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

var (
	AuthClient     AuthCx
	secret         string
	ErrInvalidType = errors.New("invalid type")
)

func init() {
	AuthClient = &http.Client{}
	secret = lib.Getenv("JWT_SECRET")
}

func Router(r chi.Router) {
	r.Handle("/github", http.HandlerFunc(VerifyToken))
}
