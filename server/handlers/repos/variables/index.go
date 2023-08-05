package variables

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type contributor struct {
	*jwt.GithubUser
	access bool
	err    error
	mut    sync.Mutex
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := jwt.GetUser(r)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Invalid repository id: %s", err.Error())
		return
	}

	// QUERY DB FOR REPO INFO
	repo := Repository{
		Repository: database.Repository{ID: repoID},
		Variables:  []database.Variable{},
	}
	err = db.getRepoByID(&repo)
	switch err {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).
			Status(http.StatusNotFound).
			Error("Repository not found")
		return

	default:
		api.NewResponse(w).ServerError(err)
		return
	}

	// CHECK FOR USER ACCESS
	wg := sync.WaitGroup{}
	contrib := newContributor(user)
	go contrib.fetchRepoAccess(repo.FullName, &wg)

	// QUERY DB FOR ENV VARIABLES
	err = db.getVariables(&repo)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err)
		return
	}

	// decrypt values
	for i := range repo.Variables {
		err := repo.Variables[i].DecryptValue()
		if err != nil {
			api.NewResponse(w).ServerError(err)
			return
		}
	}

	wg.Wait()
	access, err := contrib.result()
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadGateway).Done()
		return
	}

	if !access {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("not a contributor")
		return
	}

	api.NewResponse(w).Status(http.StatusOK).JSON(repo)
}

func newContributor(u *jwt.GithubUser) *contributor {
	return &contributor{
		GithubUser: u,
		access:     false,
		err:        nil,
		mut:        sync.Mutex{},
	}
}

// Endpoint to check for collaborators:
// https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#check-if-a-user-is-a-repository-collaborator
func (c *contributor) fetchRepoAccess(repoURL string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	c.mut.Lock()
	defer c.mut.Unlock()

	r, err := gh.New(c.Token).Get("/repos/%s/collaborators/%s", repoURL, c.Login)
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

func (c *contributor) result() (bool, error) {
	return c.access, c.err
}
