package auth

import (
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	refTok, _ := r.Cookie(api.CookieRefTok)
	refTok.Expires = time.Now().UTC().Add(-100 * time.Hour)
	accTok, _ := r.Cookie(api.CookieAccTok)
	accTok.Expires = time.Now().UTC().Add(-100 * time.Hour)
	api.NewResponse(w).
		SetCookie(refTok).
		SetCookie(accTok).
		Status(http.StatusOK)
}
