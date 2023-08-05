package repos

type githubContext interface {
	getRepo(repoName, userToken string, r *Repository) (int, error)
}

type githubClient struct{}

var (
	ghCx githubContext
)

func init() {
	ghCx = &githubClient{}
}
