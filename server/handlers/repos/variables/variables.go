package variables

import (
	"fmt"
	"net/http"

	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type variableHandler struct {
	*gorm.DB
}

var (
	Handlers *variableHandler
)

func init() {
	Handlers = &variableHandler{db.New()}
}

// Endpoint:
// https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#check-if-a-user-is-a-repository-collaborator
func getRepoAccess(c chan<- permission, repoURL string, u *jwt.GithubUser) {
	g := gh.New(u.Token)
	res, err := g.Get("/repos/%s/collaborators/%s", repoURL, u.Login)
	if err != nil {
		buf := permission{false, err}
		c <- buf
		return
	}

	switch res.StatusCode {
	case http.StatusNoContent:
		c <- permission{true, nil}
		return

	case http.StatusNotFound:
		c <- permission{false, nil}
		return

	default:
		c <- permission{false, fmt.Errorf("Github API reponses [%v]", res.Status)}
		return
	}
}
