package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api/jsonwebtoken"
	"github.com/hn275/envhub/server/api/response"
	"github.com/hn275/envhub/server/gh"
)

type Token struct {
	Code string
}

type GithubAuthToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type User struct {
	Token     string `json:"token,omitempty"`
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type JwtToken struct {
	User `json:",inline"`
	jwt.RegisteredClaims
}

func verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var token Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		response.New(w).ServerError(err)
		return
	}

	// get auth token
	result, err := githubOAuth(token.Code)
	if err != nil {
		response.New(w).ServerError(err)
		return
	}
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		response.New(w).Status(result.StatusCode).Done()
		return
	}

	var accesstoken GithubAuthToken
	if err := json.NewDecoder(result.Body).Decode(&accesstoken); err != nil {
		response.New(w).ServerError(err)
		return
	}

	// get user account
	ghResponse, err := gh.New(accesstoken.AccessToken).Get("/user")
	if err != nil {
		response.New(w).ServerError(err)
		return
	}
	defer ghResponse.Body.Close()

	var user User
	if err := json.NewDecoder(ghResponse.Body).Decode(&user); err != nil {
		response.New(w).ServerError(err)
		return
	}

	jwtToken := JwtToken{
		user,
		jwt.RegisteredClaims{
			Issuer:  "Envhub",
			Subject: user.Name,
		},
	}

	jwtStr, err := jsonwebtoken.Sign(&jwtToken)
	if err != nil {
		response.New(w).ServerError(err)
		return
	}

	response.New(w).Status(http.StatusOK).Text(jwtStr)
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

	var cx http.Client
	return cx.Do(req)
}
