package variables

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/envhubtest"
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
