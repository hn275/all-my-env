package variables

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type variableHandler struct {
	*gorm.DB
}

var (
	Handlers           *variableHandler
	errNotAContributor = errors.New("not a contributor")
	errBadGateWay      = errors.New("GitHub responded an with error")
)

func init() {
	Handlers = &variableHandler{database.New()}
}

// Endpoint:
// https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#check-if-a-user-is-a-repository-collaborator
func getRepoAccess(c chan<- error, wg *sync.WaitGroup, repoURL string, u *jwt.GithubUser) {
	wg.Add(1)
	defer wg.Done()

	res, err := gh.New(u.Token).Get("/repos/%s/collaborators/%s", repoURL, u.Login)
	if err != nil {
		c <- err
		return
	}

	switch res.StatusCode {
	case http.StatusNoContent:
		c <- nil
		return

	case http.StatusNotFound:
		c <- errNotAContributor
		return

	default:
		fmt.Fprintf(os.Stderr, "GitHub API responded with %d", res.StatusCode)
		c <- errBadGateWay
		return
	}
}
