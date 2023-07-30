package repos

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hn275/envhub/server/gh"
	"github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type repoMock struct{}

func init() {
	gh.MockClient(&repoMock{})
	jsonwebtoken.Mock()
}

func testInit() (*http.ServeMux, *httptest.ResponseRecorder) {
	r := http.NewServeMux()
	r.Handle("/test", http.HandlerFunc(Index))
	return r, &httptest.ResponseRecorder{}
}

func TestLinkRepoRecordNotFound(t *testing.T) {
	db = &mockRepoDb{findRepoErr: gorm.ErrRecordNotFound}
	s := httptest.NewServer(http.HandlerFunc(Index))
	defer s.Close()

	url := s.URL + "?show=69&page=420&sort=foo"
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	r.Header.Add("Authorization", "Bearer aslkfjsdklfj")
	c := http.Client{}
	res, err := c.Do(r)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var repos []Repository
	err = json.NewDecoder(res.Body).Decode(&repos)
	assert.Nil(t, err)
	assert.False(t, repos[0].Linked)
}

func TestLinkRepoDBError(t *testing.T) {
	db = &mockRepoDb{findRepoErr: errors.New("some error")}
	s := httptest.NewServer(http.HandlerFunc(Index))
	defer s.Close()

	url := s.URL + "?show=69&page=420&sort=foo"
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	r.Header.Add("Authorization", "Bearer aslkfjsdklfj")
	c := http.Client{}
	res, err := c.Do(r)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Empty(t, res.Body)
}

func TestLinkedRepo(t *testing.T) {
	db = &mockRepoDb{findRepoErr: nil}
	s := httptest.NewServer(http.HandlerFunc(Index))
	defer s.Close()

	url := s.URL + "?show=69&page=420&sort=foo"
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	r.Header.Add("Authorization", "Bearer aslkfjsdklfj")
	c := http.Client{}
	res, err := c.Do(r)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var repos []Repository
	err = json.NewDecoder(res.Body).Decode(&repos)
	assert.Nil(t, err)
	assert.True(t, repos[0].Linked)
}

func TestMethodAllow(t *testing.T) {
	r, w := testInit()

	methods := []string{
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/test?show=1&sort=updated&page=1", nil)
		assert.Nil(t, err)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestAllMissingAuth(t *testing.T) {
	r, w := testInit()

	req, err := http.NewRequest(http.MethodGet, "/test?show=1&sort=updated&page=1", nil)
	assert.Nil(t, err)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
}

func TestAllBadRequest(t *testing.T) {
	r, w := testInit()

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer asdlkfjfsd")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestMaxIDVal(t *testing.T) {
	mockRepos := []Repository{}
	for i := 0; i < 1000; i++ {
		id := rand.Int() // since uint64 overflows assert.Equal
		if id == math.MaxInt {
			continue
		}
		mockRepos = append(mockRepos, Repository{ID: uint64(id)})
	}
	mockRepos = append(mockRepos, Repository{ID: uint64(math.MaxInt)})
	assert.Equal(t, uint64(math.MaxInt), maxIDVal(mockRepos))
}

// MOCK
// Do implements gh.Client
func (mock *repoMock) Do(req *http.Request) (*http.Response, error) {
	buf := bytes.NewReader([]byte(mockData))
	body := ioutil.NopCloser(buf)

	res := &http.Response{
		StatusCode: 200,
		Request:    req,
		Body:       body,
	}
	return res, nil
}

const mockData = `
[
  {
    "id": 1,
    "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
    "name": "Hello-World",
    "full_name": "octocat/Hello-World"
  }
]
`
