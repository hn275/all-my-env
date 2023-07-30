package repos

import (
	"github.com/hn275/envhub/server/db"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type githubContext interface {
	getRepo(r *Repository) (int, error)
}

type githubClient struct {
	repoName string
	user     *jwt.GithubUser
}

type RepoHandler struct{ *gorm.DB }

var (
	Handlers *RepoHandler
	ghClient githubContext
)

func init() {
	Handlers = &RepoHandler{db.New()}
}
