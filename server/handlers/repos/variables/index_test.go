package variables

import (
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

	params := map[string]string{"id": "1"}
	r := envhubtest.RequestWithParam(http.MethodGet, "/", params, nil)
	r.Header.Add("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()
	Handlers.Index(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
