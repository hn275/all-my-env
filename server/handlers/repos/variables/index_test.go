package variables

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/envhubtest"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

func TestVariableIndexRequestMethods(t *testing.T) {
	mux := chi.NewMux()
	mux.Handle("/", http.HandlerFunc(Handlers.Index))
	for _, method := range envhubtest.AllowedRequestMethods(http.MethodGet) {
		r, err := http.NewRequest(method, "/", nil)
		assert.Nil(t, err)
		w := &httptest.ResponseRecorder{}
		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestVariableIndexOK(t *testing.T) {
	jwt.Mock()
	gh.MockClient(&mockGhCtxOK{})

	params := map[string]string{"repoID": "1"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()

	Handlers.Index(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var response Repository
	err := json.NewDecoder(w.Result().Body).Decode(&response)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(response.Variables))
	assert.Equal(t, "bar", response.Variables[0].Value)
}

func TestVariableIndexNotContributor(t *testing.T) {
	// mock data
	jwt.Mock()
	gh.MockClient(&mockGhCtxNotFound{})

	// test
	params := map[string]string{"repoID": "1"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()

	Handlers.Index(w, r)

	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestVariableIndexGithubError(t *testing.T) {
	jwt.Mock()
	gh.MockClient(&mockGhCtxError{})

	params := map[string]string{"repoID": "1"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()

	Handlers.Index(w, r)

	assert.Equal(t, http.StatusBadGateway, w.Result().StatusCode)
}

func TestVariableIndexNoVars(t *testing.T) {
	jwt.Mock()
	gh.MockClient(&mockGhCtxOK{})

	params := map[string]string{"repoID": "3"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	Handlers.Index(w, r)

	var b Repository
	err := json.NewDecoder(w.Result().Body).Decode(&b)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Empty(t, b.Variables)
}

func TestVariableIndexRepoNotFound(t *testing.T) {
	jwt.Mock()
	gh.MockClient(&mockGhCtxOK{})

	params := map[string]string{"repoID": "420"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	Handlers.Index(w, r)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}
