package repos

import (
	"gorm.io/gorm"
)

type githubContext interface {
	getRepo(repoName, userToken string, r *Repository) (int, error)
}

type githubClient struct {
}

type RepoHandler struct{ *gorm.DB }

var (
	ghCx githubContext
)

func init() {
	ghCx = &githubClient{}
}
