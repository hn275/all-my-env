package auth

import (
	"encoding/json"
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

	// validate access token
	c, err := getToken(r.Header.Get("Authorization"))
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// verifying jwt
	clms, err := jsonwebtoken.NewDecoder().Decode(c)
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

	// query db to get refresh token
	type dberr struct {
		err error
	}
	wg := sync.WaitGroup{}
	dbErr := dberr{}
	go func(wg *sync.WaitGroup, dbErr *dberr) {
		wg.Add(1)
		defer wg.Done()
		var ref struct {
			RefreshToken string
		}
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

	userInfo := authResponse{
		AccessToken: jwtToken,
		Name:        u.Name,
		AvatarUrl:   u.AvatarURL,
		Login:       u.Login,
	}

	wg.Wait()
	switch err := dbErr.err; err {
	case nil:
		api.NewResponse(w).Status(http.StatusOK).JSON(&userInfo)
		return

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("refresh token not found")
		return

	default:
		api.NewResponse(w).ServerError(err.Error())
		return
	}
}
