package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, http.MethodGet) {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	// validate refresh cookie
	tok, err := r.Cookie(api.CookieRefTok)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}
	if err := tok.Valid(); err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// validate access token
	jwtTok, err := getToken(r.Header.Get("Authorization"))
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// verifying jwt
	clms, err := jsonwebtoken.NewDecoder().Decode(jwtTok)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("failed to decode token: %v", err)
		return
	}

	// get user info
	userID, err := strconv.ParseUint(clms.Subject, 10, 64)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	accessToken, err := decodeAccessToken(userID, clms.AccessToken)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// query for user info
	var wg sync.WaitGroup
	var u GithubUser
	var ghErr error
	go func(wg *sync.WaitGroup, tok string, user *GithubUser, err *error) {
		wg.Add(1)
		defer wg.Done()
		res, e := gh.New(tok).Get("/user")
		defer res.Body.Close()
		if e != nil {
			*err = e
			return
		}
		if res.StatusCode != http.StatusBadRequest {
			*err = fmt.Errorf("GitHub response: %s", res.Status)
			return
		}
		*err = json.NewDecoder(res.Body).Decode(user)
	}(&wg, accessToken, &u, &ghErr)

	// query db to get refresh token
	db := database.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `SELECT refresh_token FROM users WHERE id = ? AND refresh_token = ?`
	var refreshToken string

	err = db.GetContext(ctx, &refreshToken, q, userID, tok.Value)
	if errors.Is(err, sql.ErrNoRows) {
		api.NewResponse(w).Status(http.StatusForbidden).Done()
		return
	}

	// response
	clms.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour))
	jwtToken, err := jsonwebtoken.NewEncoder().Encode(
		userID,
		clms.AccessToken,
		clms.Audience[0],
	)
	if err != nil {
		api.NewResponse(w).ServerError("%v", err)
		return
	}

	userInfo := struct {
		AccessToken string
		Name        string
		AvatarUrl   string
		Login       string
	}{
		AccessToken: jwtToken,
		Name:        u.Name,
		AvatarUrl:   u.AvatarURL,
		Login:       u.Login,
	}

	wg.Wait()
	switch ghErr {
	case nil:
		api.NewResponse(w).Status(http.StatusOK).JSON(&userInfo)
		return

	default:
		api.NewResponse(w).Status(http.StatusBadRequest).Error("GitHub API failed, try again later.")
		return
	}
}
