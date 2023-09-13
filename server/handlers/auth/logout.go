package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	refTok, err := r.Cookie(api.CookieRefTok)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	db := database.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := "UPDATE users SET refresh_token = NULL WHERE refresh_token = ?"
	_, err = db.QueryxContext(ctx, q, refTok.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.NewResponse(w).Status(http.StatusBadRequest).Error("User not found.")
			return
		}
		api.NewResponse(w).ServerError(err.Error())
	}

	refTok.Expires = time.Now().UTC().Add(-100 * time.Hour)

	api.NewResponse(w).
		SetCookie(refTok).
		Status(http.StatusOK)
}
