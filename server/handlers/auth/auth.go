package auth

import (
	"errors"
	"net/http"

	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

type AuthCx interface {
	Do(r *http.Request) (*http.Response, error)
}

type AuthHandler struct {
	*gorm.DB
}

var (
	AuthClient     AuthCx
	ErrInvalidType = errors.New("invalid type")
	Handler        *AuthHandler
)

func init() {
	AuthClient = &http.Client{}
	Handler = &AuthHandler{db.New()}
}

func Handlers() *AuthHandler {
	return &AuthHandler{db.New()}
}
