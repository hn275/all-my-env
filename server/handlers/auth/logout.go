package auth

import (
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	refTok, err := r.Cookie(api.CookieRefTok)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}
	refTok.Expires = time.Now().UTC().Add(-100 * time.Hour)

	api.NewResponse(w).
		SetCookie(refTok).
		Status(http.StatusOK)
}
