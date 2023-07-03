package gh

import (
	"net/http"

	"github.com/hn275/envhub/server/lib"
)

const (
	githubUrl = "https://api.github.com"
)

type GithubContext struct {
	token string
}

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	GithubClientID     string
	GithubClientSecret string
	GithubClient       Client
)

func init() {
	GithubClientID = lib.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret = lib.Getenv("GITHUB_CLIENT_SECRET")
	GithubClient = &http.Client{}
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

	return GithubClient.Do(req)
}
