package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/lib"
)

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var token Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("invalid credentials")
		return
	}

	// get auth token
	result, err := githubOAuth(token.Code)
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		api.NewResponse(w).Status(result.StatusCode)
		return
	}

	var accesstoken GithubAuthToken
	if err := json.NewDecoder(result.Body).Decode(&accesstoken); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	// get user account
	ghResponse, err := gh.New(accesstoken.AccessToken).Get("/user")
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	defer ghResponse.Body.Close()

	if ghResponse.StatusCode != http.StatusOK {
		api.NewResponse(w).Status(ghResponse.StatusCode)
		return
	}

	var userInfo GithubUser
	if err := json.NewDecoder(ghResponse.Body).Decode(&userInfo); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	// save user in db (if not exists)
	user := db.User{
		ID:        userInfo.ID,
		CreatedAt: lib.TimeStamp(),
		Vendor:    db.VendorGithub,
		UserName:  userInfo.Login,
	}
	if err := h.Create(&user).Error; err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	userInfo.Token = accesstoken.AccessToken
	jwtToken := JwtToken{
		userInfo,
		jwt.RegisteredClaims{
			Issuer:  "Envhub",
			Subject: userInfo.Name,
		},
	}

	jwtStr, err := Sign(&jwtToken)
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	api.NewResponse(w).Status(http.StatusOK).Text(jwtStr)
}

func githubOAuth(code string) (*http.Response, error) {
	// get auth token
	ghEndpoint := "https://github.com/login/oauth/access_token"
	v := url.Values{}
	v.Set("client_id", gh.GithubClientID)
	v.Set("client_secret", gh.GithubClientSecret)
	v.Set("code", code)
	ghEndpoint += "?" + v.Encode()

	req, err := http.NewRequest(http.MethodPost, ghEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")

	return AuthClient.Do(req)
}
