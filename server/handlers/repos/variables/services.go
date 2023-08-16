package variables

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/hn275/envhub/server/gh"
)

type contributor struct {
	userLogin string
	userTok   string
	repoURL   string
	access    bool
	err       error
	mut       sync.Mutex
	wg        sync.WaitGroup
}

// Endpoint to check for collaborators:
// https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#check-if-a-user-is-a-repository-collaborator
func (c *contributor) getRepoAccess() {
	c.wg.Add(1)
	defer c.wg.Done()
	c.mut.Lock()
	defer c.mut.Unlock()

	if c.repoURL == "" {
		c.err = errors.New("repository url not found.")
		return
	}
	if c.userLogin == "" {
		c.err = errors.New("user login not found")
		return
	}
	if c.userTok == "" {
		c.err = errors.New("user token not found")
		return
	}

	r, err := gh.New(c.userTok).Get("/repos/%s/collaborators/%s", c.repoURL, c.userLogin)
	if err != nil {
		c.err = err
		return
	}

	switch r.StatusCode {
	case http.StatusNoContent:
		c.err = nil
		c.access = true
		return

	case http.StatusNotFound:
		c.err = nil
		c.access = false
		return

	default:
		c.err = fmt.Errorf("GitHub API responded with %d\n", r.StatusCode)
		c.access = false
		return
	}
}
