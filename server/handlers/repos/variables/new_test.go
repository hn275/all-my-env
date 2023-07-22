package variables

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

type mockCtxOK struct{}
type mockGhCtxNotFound struct{}
type mockJwtToken struct{}

var mockVar = EnvVariable{"foo", "bar"}

func testInit(url string) (*httptest.ResponseRecorder, error) {
	buf, err := json.Marshal(mockVar)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(buf)

	r, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+"somejwttoken")

	m := chi.NewMux()
	m.Handle("/repos/{id}/variables/new", http.HandlerFunc(Handlers.NewVariable))

	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	return w, nil
}

func cleanup() {
	d := db.New()
	d.Where("repository_id = ? AND key = ?", 1, "foo").Delete(&db.Variable{})
}

func TestNewVariable(t *testing.T) {
	defer cleanup()

	gh.GithubClient = &mockCtxOK{}
	jwt.Decoder = &mockJwtToken{}

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	var variable db.Variable
	d := db.New()
	err = d.Where("repository_id = ? AND key = ?", 1, "foo").First(&variable).Error
	assert.Nil(t, err)
	assert.NotEqual(t, mockVar.Value, variable.Value)
}

func TestNewVariableDuplicate(t *testing.T) {
	defer cleanup()

	gh.GithubClient = &mockCtxOK{}
	jwt.Decoder = &mockJwtToken{}

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	w, err = testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
}

func TestRepoNotFound(t *testing.T) {
	gh.GithubClient = &mockCtxOK{}
	jwt.Decoder = &mockJwtToken{}

	w, err := testInit("/repos/420/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestMethodNotAllowed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(Handlers.NewVariable))
	cx := http.Client{}

	methods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPatch,
		http.MethodHead,
		http.MethodDelete,
	}
	for _, method := range methods {
		req, err := http.NewRequest(method, srv.URL, nil)
		assert.Nil(t, err)
		res, err := cx.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestNewVarNotAllowed(t *testing.T) {
	gh.GithubClient = &mockGhCtxNotFound{}
	jwt.Decoder = &mockJwtToken{}

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

// MOCK
func (m *mockCtxOK) Do(r *http.Request) (*http.Response, error) {
	res := &http.Response{
		Status:     "204 No Content",
		StatusCode: http.StatusNoContent,
		Request:    r,
	}
	return res, nil
}

func (m *mockGhCtxNotFound) Do(r *http.Request) (*http.Response, error) {
	res := &http.Response{
		Status:     "404 Not Found",
		StatusCode: http.StatusNotFound,
		Request:    r,
	}
	return res, nil
}

func (d *mockJwtToken) Decode(t string) (*jwt.JwtToken, error) {
	decoded := &jwt.JwtToken{
		GithubUser: jwt.GithubUser{Token: t, ID: 1},
	}
	return decoded, nil
}
