package auth

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm/clause"
)

// Verify token send from body, then query for user data.
// Save user in database if they don't exists
func (h *authHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
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

	// GET AUTH TOKEN
	req, err := http.NewRequest(http.MethodPost, gh.OAuthUrl(t.Code), nil)
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	req.Header.Add("accept", "application/json")

	res, err := authClient.Do(req)
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		api.NewResponse(w).Status(res.StatusCode).Done()
		return
	}

	var oauth struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(res.Body).Decode(&oauth); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	// GET USER ACCOUNT
	// save user in db if not exists
	ghResponse, err := gh.New(oauth.AccessToken).Get("/user")
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	defer ghResponse.Body.Close()

	if ghResponse.StatusCode != http.StatusOK {
		api.NewResponse(w).Status(ghResponse.StatusCode).Done()
		return
	}

	var userInfo jsonwebtoken.GithubUser
	if err := json.NewDecoder(ghResponse.Body).Decode(&userInfo); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	// save user in db (if not exists)
	user := db.User{
		ID:        userInfo.ID,
		CreatedAt: db.TimeNow(),
		Vendor:    db.VendorGithub,
		UserName:  userInfo.Login,
	}
	err = h.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	userInfo.Token = oauth.AccessToken
	jwtToken := jsonwebtoken.JwtToken{
		GithubUser: userInfo,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "EnvHub",
			Subject: userInfo.Name,
		},
	}

	jwtStr, err := jsonwebtoken.Sign(&jwtToken)
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	api.NewResponse(w).Status(http.StatusOK).Text(jwtStr)
}
