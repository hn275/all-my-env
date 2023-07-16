package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hn275/envhub/server/jsonwebtoken"
)

func GetUser(r *http.Request) (*jsonwebtoken.GithubUser, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return nil, errors.New("auth token not found")
	}

	t := strings.Split(h, " ")

	if len(t) != 2 {
		return nil, errors.New("invalid auth token")
	}

	if strings.ToLower(t[0]) != "bearer" {
		return nil, errors.New("invalid auth token")
	}

	decoded, err := jsonwebtoken.Decode(t[1])
	if err != nil {
		return nil, err
	}

	return &decoded.GithubUser, nil
}
