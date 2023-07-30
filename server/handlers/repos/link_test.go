package repos

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
)

var mockBody = database.Repository{
	FullName: "foo",
}

func TestLinkOK(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{newRepoErr: nil}
	ghCx = &mockGhCtx{code: http.StatusOK, err: nil, ownerID: 1}

	b, err := json.Marshal(mockBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Link(w, r)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	var buf bytes.Buffer
	buf.ReadFrom(w.Body)
	assert.NotEmpty(t, buf)
	assert.Equal(t, "application/text", w.Header().Get("content-type"))
}

func TestLinkDBError(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{newRepoErr: errors.New("asdf")}
	ghCx = &mockGhCtx{code: http.StatusOK, err: nil, ownerID: 1}

	b, err := json.Marshal(mockBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Link(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestLinkGithubError(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{newRepoErr: nil}
	for _, code := range []int{http.StatusForbidden, http.StatusNotFound} {
		ghCx = &mockGhCtx{code: code, err: nil, ownerID: 2}

		b, err := json.Marshal(mockBody)
		assert.Nil(t, err)
		body := ioutil.NopCloser(bytes.NewReader(b))

		r, err := http.NewRequest(http.MethodPost, "/", body)
		assert.Nil(t, err)
		r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

		w := httptest.NewRecorder()
		Link(w, r)

		assert.Equal(t, code, w.Result().StatusCode)
	}
}

func TestLinkGithubBadGateway(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{newRepoErr: nil}
	ghCx = &mockGhCtx{code: http.StatusMovedPermanently, err: nil, ownerID: 2}

	b, err := json.Marshal(mockBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Link(w, r)

	assert.Equal(t, http.StatusBadGateway, w.Result().StatusCode)
}
