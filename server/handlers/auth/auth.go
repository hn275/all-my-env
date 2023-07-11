package auth

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/lib"
	"gorm.io/gorm"
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

type AuthHandler struct {
	*gorm.DB
}

var (
	AuthClient     AuthCx
	secret         string
	ErrInvalidType = errors.New("invalid type")
	Handler        *AuthHandler
)

func init() {
	AuthClient = &http.Client{}
	secret = lib.Getenv("JWT_SECRET")
	Handler = &AuthHandler{db.New()}
}

func Router(r chi.Router) {
	r.Handle("/github", http.HandlerFunc(Handler.VerifyToken))
}
