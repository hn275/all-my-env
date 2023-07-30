package repos

type mockGhCtx struct {
	code    int
	err     error
	ownerID uint64
}

func (ctx *mockGhCtx) getRepo(repoName, userToken string, repo *Repository) (int, error) {
	repo.ID = 1
	repo.FullName = "foo"
	repo.Owner.ID = ctx.ownerID
	return ctx.code, ctx.err
}
