package gh

import (
	"net/http"
)

const (
	githubUrl = "https://api.github.com"
)

var (
	GithubClientID     string
	GithubClientSecret string
)

type GithubContext struct {
	token string
}

type GithubClient struct {
	http.Client
}

func init() {
	GithubClientID = lib.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret = lib.Getenv("GITHUB_CLIENT_SECRET")
}

func New(token string) *GithubContext {
	return &GithubContext{token}
}

func (g *GithubContext) Get(path string) (*http.Response, error) {
	if path[0] != '/' {
		path = "/" + path
	}
	url := githubUrl + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Envhub")
	req.Header.Add("Authorization", "Bearer "+g.token)
	req.Header.Add("Accent", "application/vnd.github+json")

	var cx http.Client
	return cx.Do(req)
}
