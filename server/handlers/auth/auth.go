package auth

import (
	"net/http"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type authCtx interface {
	Do(r *http.Request) (*http.Response, error)
}

type authHandler struct {
	*gorm.DB
}

var (
	authClient authCtx
	Handler    *authHandler
)

func init() {
	authClient = &http.Client{}
	Handler = &authHandler{database.New()}
}
