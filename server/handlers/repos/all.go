package repos

import (
	"encoding/json"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
)

func (h *RepoHandler) All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).
			Status(http.StatusMethodNotAllowed).
			Error("http method not allowed")
		return
	}

	user, err := api.GetUser(r)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	// GET REQUEST QUERY PARAM
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")
	show := r.URL.Query().Get("show")
	if page == "" || sort == "" || show == "" {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("missing required queries")
		return
	}
	params := map[string]string{
		"page":     page,
		"sort":     sort,
		"per_page": show,
	}

	// NOTE: Since we are only interested in the repo that got sent back
	// by Github, this ops won't be a go routine.
	// get repos from github, then query db for the id of the same set of repos.

	// GET REPOS FROM GITHUB
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-the-authenticated-user
	var repos []Repository
	ghCtx := gh.New(user.Token).Params(params)

	res, err := ghCtx.Get("/user/repos")
	defer res.Body.Close()
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	ids := make([]uint, len(repos))
	for i, repo := range repos {
		ids[i] = repo.ID
	}

	// GET REPO ID's FROM DB
	var dbRepos []uint
	trx := h.Table(db.TableRepos)
	trx.Select("id")
	trx.Where("user_id = ? AND id IN ?", user.ID, ids[:])
	if err := trx.Find(&dbRepos).Error; err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

repoLoop:
	for i, repo := range repos {
		for _, id := range dbRepos {
			if id == repo.ID {
				repos[i].Linked = true
				continue repoLoop
			}
		}
	}

	api.NewResponse(w).
		Header("Cache-Control", "max-age=30").
		Status(http.StatusOK).
		JSON(&repos)
}
