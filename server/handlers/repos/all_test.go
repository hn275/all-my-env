package repos_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/handlers/repos"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

type repoMock struct{}
type jwtMock struct{}

func init() {
	gh.GithubClient = &repoMock{}
	jsonwebtoken.Decoder = &jwtMock{}
}

func testInit() (*chi.Mux, *httptest.ResponseRecorder) {
	r := chi.NewMux()
	r.Handle("/test", http.HandlerFunc(repos.Handlers.All))
	return r, &httptest.ResponseRecorder{}
}

func TestLinkedRepo(t *testing.T) {
	var ghRepos []repos.Repository
	if err := json.Unmarshal([]byte(mockData), &ghRepos); err != nil {
		panic(err)
	}

	mux, w := testInit()
	r, err := http.NewRequest(http.MethodGet, "/test?show=69&page=420&sort=foo", nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("Authorization", "Bearer "+jwtToken)

	mux.ServeHTTP(w, r)
	result := w.Result()
	defer result.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestMethodAllow(t *testing.T) {
	r, w := testInit()

	methods := []string{
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/test?show=1&sort=updated&page=1", nil)
		assert.Nil(t, err)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestAllMissingAuth(t *testing.T) {
	r, w := testInit()

	req, err := http.NewRequest(http.MethodGet, "/test?show=1&sort=updated&page=1", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestAllBadRequest(t *testing.T) {
	r, w := testInit()

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+jwtToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

// MOCK
// Do implements gh.Client
func (mock *repoMock) Do(req *http.Request) (*http.Response, error) {
	buf := bytes.NewReader([]byte(mockData))
	body := ioutil.NopCloser(buf)

	res := &http.Response{
		StatusCode: 200,
		Request:    req,
		Body:       body,
	}
	return res, nil
}

// Decode implements jsonwebtoken.JsonWebToken.
func (*jwtMock) Decode(_ string) (*jsonwebtoken.JwtToken, error) {
	user := jsonwebtoken.GithubUser{
		Token:     "asdf",
		ID:        123,
		Login:     "foo",
		AvatarUrl: "foobar.com",
		Name:      "foo",
		Email:     "foo@bar.com",
	}
	t := &jsonwebtoken.JwtToken{
		GithubUser: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "EnvHub",
			Subject:   "foo",
			Audience:  []string{},
			ExpiresAt: &jwt.NumericDate{},
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  &jwt.NumericDate{},
			ID:        "123",
		},
	}
	return t, nil
}

const jwtToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6Imdob19MemhsbXNkU3g3b3phdUs0ZFFKejcyMmRkOFJ6bWo0SloxMzkiLCJpZCI6OTcxNDM1OTYsImxvZ2luIjoiaG4yNzUiLCJhdmF0YXJfdXJsIjoiaHR0cHM6Ly9hdmF0YXJzLmdpdGh1YnVzZXJjb250ZW50LmNvbS91Lzk3MTQzNTk2P3Y9NCIsIm5hbWUiOiJIYWwiLCJlbWFpbCI6ImhhbG5fMDFAcHJvdG9uLm1lIiwiaXNzIjoiRW52aHViIiwic3ViIjoiSGFsIn0.-tMfdpMMnxmvM-oMSyhtw8_QzrJ8AWwUNUEzOCQGh4Y`
const mockData = `
[
  {
    "id": 123,
    "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
    "name": "Hello-World",
    "full_name": "octocat/Hello-World",
    "owner": {
      "login": "octocat",
      "id": 1,
      "node_id": "MDQ6VXNlcjE=",
      "avatar_url": "https://github.com/images/error/octocat_happy.gif",
      "gravatar_id": "",
      "url": "https://api.github.com/users/octocat",
      "html_url": "https://github.com/octocat",
      "followers_url": "https://api.github.com/users/octocat/followers",
      "following_url": "https://api.github.com/users/octocat/following{/other_user}",
      "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
      "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
      "organizations_url": "https://api.github.com/users/octocat/orgs",
      "repos_url": "https://api.github.com/users/octocat/repos",
      "events_url": "https://api.github.com/users/octocat/events{/privacy}",
      "received_events_url": "https://api.github.com/users/octocat/received_events",
      "type": "User",
      "site_admin": false
    },
    "private": false,
    "html_url": "https://github.com/octocat/Hello-World",
    "description": "This your first repo!",
    "fork": false,
    "url": "https://api.github.com/repos/octocat/Hello-World",
    "archive_url": "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
    "assignees_url": "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
    "blobs_url": "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
    "branches_url": "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
    "collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
    "comments_url": "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
    "commits_url": "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
    "compare_url": "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
    "contents_url": "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
    "contributors_url": "https://api.github.com/repos/octocat/Hello-World/contributors",
    "deployments_url": "https://api.github.com/repos/octocat/Hello-World/deployments",
    "downloads_url": "https://api.github.com/repos/octocat/Hello-World/downloads",
    "events_url": "https://api.github.com/repos/octocat/Hello-World/events",
    "forks_url": "https://api.github.com/repos/octocat/Hello-World/forks",
    "git_commits_url": "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
    "git_refs_url": "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
    "git_tags_url": "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
    "git_url": "git:github.com/octocat/Hello-World.git",
    "issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
    "issue_events_url": "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
    "issues_url": "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
    "keys_url": "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
    "labels_url": "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
    "languages_url": "https://api.github.com/repos/octocat/Hello-World/languages",
    "merges_url": "https://api.github.com/repos/octocat/Hello-World/merges",
    "milestones_url": "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
    "pulls_url": "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
    "releases_url": "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
    "ssh_url": "git@github.com:octocat/Hello-World.git",
    "stargazers_url": "https://api.github.com/repos/octocat/Hello-World/stargazers",
    "statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
    "subscribers_url": "https://api.github.com/repos/octocat/Hello-World/subscribers",
    "subscription_url": "https://api.github.com/repos/octocat/Hello-World/subscription",
    "tags_url": "https://api.github.com/repos/octocat/Hello-World/tags",
    "teams_url": "https://api.github.com/repos/octocat/Hello-World/teams",
    "trees_url": "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
    "clone_url": "https://github.com/octocat/Hello-World.git",
    "mirror_url": "git:git.example.com/octocat/Hello-World",
    "hooks_url": "https://api.github.com/repos/octocat/Hello-World/hooks",
    "svn_url": "https://svn.github.com/octocat/Hello-World",
    "homepage": "https://github.com",
    "language": null,
    "forks_count": 9,
    "stargazers_count": 80,
    "watchers_count": 80,
    "size": 108,
    "default_branch": "master",
    "open_issues_count": 0,
    "is_template": true,
    "topics": ["octocat", "atom", "electron", "api"],
    "has_issues": true,
    "has_projects": true,
    "has_wiki": true,
    "has_pages": false,
    "has_downloads": true,
    "archived": false,
    "disabled": false,
    "visibility": "public",
    "pushed_at": "2011-01-26T19:06:43Z",
    "created_at": "2011-01-26T19:01:12Z",
    "updated_at": "2011-01-26T19:14:43Z",
    "permissions": {
      "admin": false,
      "push": false,
      "pull": true
    },
    "allow_rebase_merge": true,
    "template_repository": null,
    "temp_clone_token": "ABTLWHOULUVAXGTRYU7OC2876QJ2O",
    "allow_squash_merge": true,
    "allow_auto_merge": false,
    "delete_branch_on_merge": true,
    "allow_merge_commit": true,
    "subscribers_count": 42,
    "network_count": 0,
    "license": {
      "key": "mit",
      "name": "MIT License",
      "url": "https://api.github.com/licenses/mit",
      "spdx_id": "MIT",
      "node_id": "MDc6TGljZW5zZW1pdA==",
      "html_url": "https://github.com/licenses/mit"
    },
    "forks": 1,
    "open_issues": 1,
    "watchers": 1
  }
]
`
