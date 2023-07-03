package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/gh"
)

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var token Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		api.NewResponse(w).ServerError(err)
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

	var user User
	if err := json.NewDecoder(ghResponse.Body).Decode(&user); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	user.Token = accesstoken.AccessToken
	jwtToken := JwtToken{
		user,
		jwt.RegisteredClaims{
			Issuer:  "Envhub",
			Subject: user.Name,
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

	req, err := http.NewRequest("POST", ghEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")

	return AuthClient.Do(req)
}
