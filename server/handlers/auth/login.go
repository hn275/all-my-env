package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
)

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
		Login:        u.Login,
		Email:        u.Email,
		RefreshToken: sql.NullString{String: refreshToken, Valid: true},
	}

	db := database.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `
	INSERT INTO users (id,login,email,refresh_token) 
	VALUES (:id,:login,:email,:refresh_token)
	ON DUPLICATE KEY UPDATE 
		login = :login, 
		email = :email,
		refresh_token = :refresh_token;`
	_, err = db.NamedExecContext(ctx, q, user)
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
	// TODO: move this type to auth.go or something
	userInfo := UserAuthData{
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
