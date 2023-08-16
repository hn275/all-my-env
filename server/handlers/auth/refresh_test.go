package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/envhubtest"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenMethodNotAllowed(t *testing.T) {
	methods := envhubtest.AllowedRequestMethods(http.MethodGet)
	for _, method := range methods {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(method, "/", nil)
		assert.Nil(t, err)
		RefreshToken(w, r)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestRefreshTokenAuthNotFound(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	c := http.Cookie{
		Name:    api.CookieRefTok,
		Value:   "foobarbaz",
		Expires: time.Now().Add(100 * time.Hour),
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()
	RefreshToken(w, r)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestRefreshTokenCookieNotFound(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer sometoken")
	w := httptest.NewRecorder()
	RefreshToken(w, r)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}
