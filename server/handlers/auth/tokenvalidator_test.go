package auth

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/envhubtest"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

func TestTokenValidatorNoCookie(t *testing.T) {
	jsonwebtoken.MockDecoder(&mockJwt{})

	m := chi.NewMux()
	m.Use(TokenValidator)
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := envhubtest.RequestWithParam(http.MethodGet, "/", map[string]string{}, nil)
	r.Header.Add("Authorization", "Bearer token")

	m.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestTokenValidatorFailedJWT(t *testing.T) {
	jsonwebtoken.MockDecoder(&mockJwtFailed{})

	m := chi.NewMux()
	m.Use(TokenValidator)
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := envhubtest.RequestWithParam(http.MethodGet, "/", map[string]string{}, nil)
	cookie := http.Cookie{
		Name:     api.CookieRefTok,
		Value:    "sometoken",
		Path:     "",
		Domain:   "",
		Expires:  time.Now().Add(100 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	r.AddCookie(&cookie)
	r.Header.Add("Authorization", "Bearer token")
	buf := bytes.Buffer{}
	buf.ReadFrom(w.Body)

	m.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestTokenValidatorOK(t *testing.T) {
	jsonwebtoken.MockDecoder(&mockJwt{})

	m := chi.NewMux()
	m.Use(TokenValidator)
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	w := httptest.NewRecorder()
	r := envhubtest.RequestWithParam(http.MethodGet, "/", map[string]string{}, nil)
	cookie := http.Cookie{
		Name:     api.CookieRefTok,
		Value:    "sometoken",
		Path:     "",
		Domain:   "",
		Expires:  time.Now().Add(100 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	r.AddCookie(&cookie)
	r.Header.Add("Authorization", "Bearer token")
	buf := bytes.Buffer{}
	buf.ReadFrom(w.Body)

	m.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

type mockJwt struct{}

func (m *mockJwt) Decode(_ string) (*jsonwebtoken.AuthClaim, error) {
	accessToken, _ := encodeAccessToken(uint64(1), "someaccesstoken")

	auth := jsonwebtoken.AuthClaim{
		AccessToken: accessToken,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "EnvHub",
			Subject:   "1",
			Audience:  []string{"Foo"},
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(10000 * time.Hour)},
			NotBefore: &jwt.NumericDate{Time: time.Now().UTC()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		},
	}
	return &auth, nil
}

type mockJwtFailed struct{}

func (m *mockJwtFailed) Decode(string) (*jsonwebtoken.AuthClaim, error) {
	return nil, errors.New("some error")
}
