package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm/clause"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	AvatarUrl   string `json:"avatar_url"`
	Login       string `json:"login"`
}

// Verify token send from body, then query for user data.
// Save user in database if they don't exists
func LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var t struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("invalid credentials")
		return
	}

	// GET GITHUB ACCESS TOKEN
	req, err := http.NewRequest(http.MethodPost, gh.OAuthUrl(t.Code), nil)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	req.Header.Add("accept", "application/json")

	res, err := authClient.Do(req)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(res)
		return
	}

	var oauth struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(res.Body).Decode(&oauth); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// GET USER ACCOUNT
	// save user in db if not exists
	ghRes, err := gh.New(oauth.AccessToken).Get("/user")
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	defer ghRes.Body.Close()

	if ghRes.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(ghRes)
		return
	}

	var u GithubUser
	if err := json.NewDecoder(ghRes.Body).Decode(&u); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// GET REFRESH TOKEN
	refreshToken, err := refreshToken(u.ID)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// SAVE TO DB
	user := database.User{
		ID:           u.ID,
		LastLogin:    database.TimeNow(),
		Login:        u.Login,
		Email:        u.Email,
		RefreshToken: refreshToken,
	}

	db := database.NewGorm()
	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// NEW JWT
	maskedAccessTok, err := encodeAccessToken(u.ID, oauth.AccessToken)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	accessJWT, err := jsonwebtoken.NewEncoder().Encode(u.ID, maskedAccessTok, u.Login)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// SET RESPONSES
	userInfo := authResponse{
		AccessToken: accessJWT,
		Name:        u.Name,
		AvatarUrl:   u.AvatarURL,
		Login:       u.Login,
	}
	refreshCookie := http.Cookie{
		Name:     api.CookieRefTok,
		Value:    refreshToken,
		Expires:  time.Now().UTC().Add(365 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	api.NewResponse(w).
		SetCookie(&refreshCookie).
		Status(http.StatusOK).
		JSON(&userInfo)
}
