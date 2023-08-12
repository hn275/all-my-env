package auth

import (
	"encoding/json"
	"errors"
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
	"gorm.io/gorm"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, http.MethodGet) {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	// validate auth token
	authToken, err := validateAuthToken(r.Header.Get("Authorization"))
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// verifying jwt
	clms, err := jsonwebtoken.NewDecoder().Decode(authToken)
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

	// validate ref token
	tok, err := r.Cookie(api.CookieRefTok)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}
	if err := tok.Valid(); err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	wg := sync.WaitGroup{}
	type DbErr struct {
		err error
	}
	dbErr := DbErr{nil}
	go func(wg *sync.WaitGroup, dbErr *DbErr) {
		wg.Add(1)
		defer wg.Done()
		var ref struct{ RefreshToken string }
		dbErr.err = database.New().
			Table(database.TableUsers).
			Select("refresh_token").
			Where("id = ? AND refresh_token = ?", userID, tok.Value).
			First(&ref).Error
	}(&wg, &dbErr)

	// get user info
	res, err := gh.New(string(accessToken)).Get("/user")
	if err != nil {
		api.NewResponse(w).Error(err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(res)
		return
	}

	var u GithubUser
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	clms.ExpiresAt = jwt.NewNumericDate(time.Now().UTC().Add(7 * 24 * time.Hour))
	jwtToken, err := jsonwebtoken.NewEncoder().Encode(userID, clms.AccessToken, clms.Audience[0])
	if err != nil {
		api.NewResponse(w).ServerError("%v", err)
		return
	}

	userInfo := authResponse{
		AccessToken: jwtToken,
		Name:        u.Name,
		AvatarUrl:   u.AvatarURL,
		Login:       u.Login,
	}

	wg.Wait()
	if dbErr.err != nil {
		if errors.Is(dbErr.err, gorm.ErrRecordNotFound) {
			api.NewResponse(w).Status(http.StatusForbidden).Error("user not found")
			return
		}
		api.NewResponse(w).ServerError(dbErr.err.Error())
		return
	}
	api.NewResponse(w).Status(http.StatusOK).JSON(&userInfo)
}
