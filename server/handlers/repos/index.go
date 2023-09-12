package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

	// GET REPO ID's FROM DB
	// NOTE: since `repo.ID` is from github and is only of type uint32
	// string interpolation is fine, even though is bad practice, it is what it is :(
	ids := make([]string, len(repos))
	for i, repo := range repos {
		ids[i] = fmt.Sprintf("%d", repo.ID)
	}

	arr := fmt.Sprintf("%s", strings.Join(ids, ","))
	q := `SELECT id FROM repositories WHERE id IN (` + arr + `);`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.New()
	rows, err := db.QueryContext(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	repoMap := make(map[uint32]interface{})
	for rows.Next() {
		var i uint32
		if err := rows.Scan(&i); err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		repoMap[i] = struct{}{}
	}

	for i, repo := range repos {
		_, ok := repoMap[repo.ID]
		repos[i].Linked = ok
		repos[i].IsOwner = repo.Owner.ID == user.ID
	}

	api.NewResponse(w).
		Header("Cache-Control", "max-age=15").
		Status(http.StatusOK).
		JSON(&repos)
}
