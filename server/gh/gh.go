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
	param map[string]string
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
	return &GithubContext{token, nil}
}

func (ctx *GithubContext) Params(p map[string]string) *GithubContext {
	ctx.param = p
	return ctx
}

func (ctx *GithubContext) Get(path string) (*http.Response, error) {
	if path[0] != '/' {
		path = "/" + path
	}
	url := githubUrl + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if ctx.param != nil {
		q := req.URL.Query()
		for k, v := range ctx.param {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("User-Agent", "Envhub")
	req.Header.Add("Authorization", "Bearer "+ctx.token)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	return GithubClient.Do(req)
}
