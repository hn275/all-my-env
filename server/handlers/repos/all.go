package repos

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	"github.com/hn275/envhub/server/gh"
)

type ghRequest struct {
	data []Repository
	code int
}

var errSessionExpired error = errors.New("session expired")

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

	// GET REPOS FROM GITHUB
	// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repositories-for-the-authenticated-user
	ghChan := make(chan ghRequest, 1)
	defer close(ghChan)
	ghCtx := gh.New(user.Token).Params(params)
	go func(ghChan chan<- ghRequest, ghCtx *gh.GithubContext) {
		var repos []Repository
		ghResult := ghRequest{}

		res, err := ghCtx.Get("/user/repos")
		defer res.Body.Close()
		if err != nil {
			ghResult.code = http.StatusInternalServerError
			ghResult.data = nil
			ghChan <- ghResult
			return
		}

		if res.StatusCode != http.StatusOK {
			ghResult.code = res.StatusCode
			ghResult.data = nil
			ghChan <- ghResult
			return
		}

		if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
			ghResult.code = http.StatusInternalServerError
			ghResult.data = nil
			ghChan <- ghResult
			return
		}

		ghResult.code = res.StatusCode
		ghResult.data = repos
		ghChan <- ghResult
	}(ghChan, ghCtx)

	// GET REPOS FROM DB
	var dbRepos []db.Repository
	if err := h.Where("user_id = ?", user.ID).Find(&dbRepos).Error; err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	// SERIALIZE
	for {
		select {
		case m, ok := <-ghChan:
			if !ok {
				err := errors.New("github context channel closed")
				api.NewResponse(w).ServerError(err)
				return
			}

			if m.code != http.StatusOK {
				api.NewResponse(w).Status(m.code)
				return
			}

		outer:
			for _, repo := range dbRepos {
				for i, r := range m.data {
					if repo.ID == r.ID {
						m.data[i].Linked = true
						continue outer
					}
				}
			}
			api.NewResponse(w).Status(200).JSON(&m.data)
			return
		default:
			continue
		}
	}
}
