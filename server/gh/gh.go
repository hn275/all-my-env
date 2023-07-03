package gh

import (
	"fmt"
	"net/http"
	"os"
)

const (
	githubUrl = "https://api.github.com"
)

var (
	GithubClientID     string
	GithubClientSecret string
)

type GithubClient struct {
	token string
}

func init() {
	GithubClientID = os.Getenv("GITHUB_CLIENT_ID")
	if GithubClientID == "" {
		panic("GITHUB_CLIENT_ID not set")
	}
	GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if GithubClientSecret == "" {
		panic("GITHUB_CLIENT_SECRET not set")
	}
}

func New(token string) *GithubClient {
	return &GithubClient{token}
}

func (g *GithubClient) Get(path string) (*http.Response, error) {
	if path[0] != '/' {
		path = "/" + path
	}
	url := fmt.Sprintf("%s/%s", githubUrl, path)

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
