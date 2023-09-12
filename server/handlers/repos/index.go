package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).
			Status(http.StatusMethodNotAllowed).
			Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	// GET REQUEST QUERY PARAM
	params := make(map[string]string)
	params["page"] = r.URL.Query().Get("page")
	params["sort"] = r.URL.Query().Get("sort")
	params["show"] = r.URL.Query().Get("show")
	for _, v := range []string{"page", "sort", "show"} {
		if c := params[v]; c == "" {
			api.NewResponse(w).Status(http.StatusBadRequest).Error("missing required query: %s", c)
			return
		}
	}

	// GET REPOS FROM GITHUB
	// NOTE: Since we are only interested in the repo that got sent back
	// by Github, this ops won't be a go routine.
	// get repos from github, then query db for the id of the same set of repos.
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-the-authenticated-user
	var repos []Repository
	repoRes, err := gh.New(user.Token).Params(params).Get("/user/repos")
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	defer repoRes.Body.Close()

	if repoRes.StatusCode != http.StatusOK {
		api.NewResponse(w).ForwardBadRequest(repoRes)
		return
	}

	if err := json.NewDecoder(repoRes.Body).Decode(&repos); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	ids := make([]uint32, len(repos))
	for i, repo := range repos {
		ids[i] = repo.ID
		repos[i].IsOwner = repo.Owner.ID == user.ID
	}

	// GET REPO ID's FROM DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	q := `SELECT id FROM repositories WHERE user_id = ? AND id IN (?)`
	var savedReposIDs []uint32
	db := database.New()
	err = db.GetContext(ctx, savedReposIDs, q, user.ID, ids)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	repoMap := make(map[uint32]interface{})
	for i := range repos {
		repoMap[repos[i].ID] = struct{}{}
	}

	for i := range repos {
		_, ok := repoMap[repos[i].ID]
		repos[i].Linked = ok
	}

	api.NewResponse(w).
		Header("Cache-Control", "max-age=15").
		Status(http.StatusOK).
		JSON(&repos)
}
