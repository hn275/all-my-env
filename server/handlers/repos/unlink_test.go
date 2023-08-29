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

var mockUnlinkBody = database.Repository{
	FullName: "foo",
}

func TestUnlinkOK(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{deleteRepoErr: nil}
	ghCx = &mockGhCtx{code: http.StatusOK, err: nil, ownerID: 1}

	b, err := json.Marshal(mockUnlinkBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodDelete, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Unlink(w, r)

	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
}

func TestUnlinkDBError(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{deleteRepoErr: errors.New("asdf")}
	ghCx = &mockGhCtx{code: http.StatusOK, err: nil, ownerID: 1}

	b, err := json.Marshal(mockUnlinkBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodDelete, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Unlink(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestUnlinkGithubError(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{deleteRepoErr: nil}
	for _, code := range []int{http.StatusForbidden, http.StatusNotFound} {
		ghCx = &mockGhCtx{code: code, err: nil, ownerID: 2}

		b, err := json.Marshal(mockUnlinkBody)
		assert.Nil(t, err)
		body := ioutil.NopCloser(bytes.NewReader(b))

		r, err := http.NewRequest(http.MethodDelete, "/", body)
		assert.Nil(t, err)
		r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

		w := httptest.NewRecorder()
		Unlink(w, r)

		assert.Equal(t, code, w.Result().StatusCode)
	}
}

func TestUnlinkGithubBadGateway(t *testing.T) {
	jsonwebtoken.Mock()
	db = &mockRepoDb{deleteRepoErr: nil}
	ghCx = &mockGhCtx{code: http.StatusMovedPermanently, err: nil, ownerID: 2}

	b, err := json.Marshal(mockUnlinkBody)
	assert.Nil(t, err)
	body := ioutil.NopCloser(bytes.NewReader(b))

	r, err := http.NewRequest(http.MethodDelete, "/", body)
	assert.Nil(t, err)
	r.Header.Add("Authorization", "Bearer asdklfjsdklfj")

	w := httptest.NewRecorder()
	Unlink(w, r)

	assert.Equal(t, http.StatusBadGateway, w.Result().StatusCode)
}
