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

type mockGhCtxOK struct{}
type mockGhCtxNotFound struct{}
type mockGhCtxError struct{}

var mockVar = db.Variable{
	Key:   "test_foo",
	Value: "test_bar",
}

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
	m.Handle("/repos/{repoID}/variables/new", http.HandlerFunc(Handlers.NewVariable))

	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	return w, nil
}

func cleanup() {
	d := db.New()
	d.Where("key = ?", "test_foo").Delete(&db.Variable{})
}

func TestNewVariable(t *testing.T) {
	defer cleanup()

	gh.MockClient(&mockGhCtxOK{})
	jwt.Mock()

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

	gh.MockClient(&mockGhCtxOK{})
	jwt.Mock()

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	w, err = testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
}

func TestWriteAccess(t *testing.T) {
	// testing no repo not found
	gh.MockClient(&mockGhCtxOK{})
	jwt.Mock()

	w, err := testInit("/repos/420/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	// testing write-access ok
	defer cleanup()
	w, err = testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	var variable db.Variable
	d := db.New()
	err = d.Where("repository_id = ? AND key = ?", 1, "foo").First(&variable).Error
	assert.Nil(t, err)
	assert.NotEqual(t, mockVar.Value, variable.Value)
}

func TestInvalidRepoID(t *testing.T) {
	gh.MockClient(&mockGhCtxOK{})
	jwt.Mock()

	w, err := testInit("/repos/lkasjdf/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
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

func TestGithubServerNotFound(t *testing.T) {
	gh.MockClient(&mockGhCtxNotFound{})
	jwt.Mock()

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestGithubServerError(t *testing.T) {
	gh.GithubClient = &mockGhCtxError{}
	jwt.Mock()

	w, err := testInit("/repos/1/variables/new")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadGateway, w.Result().StatusCode)
}

// MOCK
func (m *mockGhCtxOK) Do(r *http.Request) (*http.Response, error) {
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

func (m *mockGhCtxError) Do(r *http.Request) (*http.Response, error) {
	res := &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: http.StatusInternalServerError,
		Request:    r,
	}
	return res, nil
}
