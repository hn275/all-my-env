package repos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).
			Status(http.StatusMethodNotAllowed).
			Error("http method not allowed")
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
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
	res, err := gh.New(user.Token).Params(params).Get("/user/repos")
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(res)
		return
	}

	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	ids := make([]uint64, len(repos))
	for i, repo := range repos {
		ids[i] = repo.ID
		repos[i].IsOwner = repo.Owner.ID == user.ID
	}

	// GET REPO ID's FROM DB
	dbRepos, err := db.findRepo(user.ID, ids[:])
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	repoMap := repoMap(dbRepos)
	for i := range repos {
		counter, ok := repoMap[repos[i].ID]
		if !ok {
			continue
		}
		repos[i].Linked = true
		repos[i].VariableCounter = counter
	}
	api.NewResponse(w).
		Header("Cache-Control", "max-age=30").
		Status(http.StatusOK).
		JSON(&repos)
}

func maxIDVal(ids []Repository) uint64 {
	var max uint64 = 0
	for _, v := range ids {
		if v.ID > max {
			max = v.ID
		}
	}

	return max
}

func repoMap(r []database.Repository) map[uint64]uint8 {
	m := make(map[uint64]uint8)
	for _, v := range r {
		m[v.ID] = uint8(v.VariableCount)
	}
	return m
}
