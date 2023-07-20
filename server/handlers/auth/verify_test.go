package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/handlers/auth"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

const (
	test_token = "sometesttoken"
)

type githubMock struct{}
type githubMockFailed struct{}
type authCxMock struct{}

func init() {
	gh.GithubClient = &githubMock{}
	auth.AuthClient = &authCxMock{}
}

func testInit() (*chi.Mux, *bytes.Reader) {
	m := chi.NewMux()
	m.Handle("/auth/github", http.HandlerFunc(auth.Handler.VerifyToken))

	token := struct {
		Code string `json:"code"`
	}{
		Code: "sometestcode",
	}
	b, _ := json.Marshal(token)
	body := bytes.NewReader(b)

	return m, body

}

func TestVerifyTokenMethodNotAllowed(t *testing.T) {
	m, body := testInit()

	var w httptest.ResponseRecorder
	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/auth/github", body)
		assert.Nil(t, err)
		m.ServeHTTP(&w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestVerifyTokenOK(t *testing.T) {
	m, body := testInit()

	w := httptest.NewRecorder()
	post, _ := http.NewRequest(http.MethodPost, "/auth/github", body)

	m.ServeHTTP(w, post)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	contentType := w.Header().Get("content-type")
	assert.Equal(t, "application/text", contentType)

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(w.Result().Body)
	assert.Nil(t, err)
	assert.NotEmpty(t, buf.Bytes())

	token, err := jsonwebtoken.NewDecoder().Decode(buf.String())
	assert.Nil(t, err)
	assert.NotEmpty(t, token.Token)
	assert.NotEmpty(t, token.Email)
	assert.NotEmpty(t, token.Name)
	assert.NotEmpty(t, token.Login)
	assert.NotEmpty(t, token.GithubUser.ID)
	assert.Equal(t, token.Issuer, "Envhub")
	assert.Equal(t, token.Token, test_token)
}

// implement gh.Client
func (m *githubMock) Do(req *http.Request) (*http.Response, error) {
	user := jsonwebtoken.GithubUser{
		ID:        123,
		Login:     "hn275",
		AvatarUrl: "https://avatars.githubusercontent.com/u/97143596?v=4",
		Name:      "Hal",
		Email:     "email@email.com",
	}
	b, _ := json.Marshal(&user)
	buf := bytes.NewReader(b)
	body := ioutil.NopCloser(buf)

	res := http.Response{
		StatusCode: 200,
		Request:    req,
		Body:       body,
	}

	return &res, nil
}

// implement auth.AuthCx
func (m *authCxMock) Do(req *http.Request) (*http.Response, error) {
	data := fmt.Sprintf(
		`{"access_token":"%s","scope":"repo,gist","token_type":"bearer"}`,
		test_token,
	)
	buf := bytes.NewReader([]byte(data))
	body := ioutil.NopCloser(buf)

	res := http.Response{
		StatusCode: 200,
		Body:       body,
		Request:    req,
	}

	return &res, nil
}

// implement gh.Client
func (m *githubMockFailed) Do(r *http.Request) (*http.Response, error) {}
