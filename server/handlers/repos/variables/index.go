package variables

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/hn275/envhub/server/gh"
	"gorm.io/gorm"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	rCtx := api.NewContext(r)
	user, err := rCtx.User()
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	// get repo info
	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("failed to parse repository id: %s", err.Error())
		return
	}

	repo := database.Repository{ID: uint32(repoID)}
	err = db.Find(&repo).Error
	switch err {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).Status(http.StatusNotFound).Error("repository not found")
		return

	default:
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	env := make([]database.Variable, 0)

	// err = db.getVariables(&env, repo.ID)
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	api.NewResponse(w).ServerError(err.Error())
	// 	return
	// }
	for i := range env {
		err = env[i].DecryptValue()
		if err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
	}

	// GET CONTRIBUTOR LIST
	var contributors []contributor
	var wg sync.WaitGroup
	var apiError error
	go getContributors(&wg, user.Token, repo.FullName, &contributors, apiError)

	// get all contributors
	var u []struct {
		UserID int64
	}

	// get contributors with write access
	err = db.Table(database.TablePermissions).
		Select("user_id").
		Where("repository_id = ?", repoID).Find(&u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	accessMap := make(map[int64]interface{})
	for _, user := range u {
		accessMap[user.UserID] = struct{}{}
	}

	wg.Wait()
	// check for repo access
	hasAccess := false
	for _, u := range u {
		if u.UserID == int64(user.ID) {
			hasAccess = true
			break
		}
	}
	if !hasAccess {
		api.NewResponse(w).Status(http.StatusForbidden).Error("not a contributor.")
		return
	}

	// combine data from gh api and database
	for i, u := range contributors {
		_, ok := accessMap[int64(u.ID)]
		contributors[i].WriteAccess = ok
	}

	response := map[string]any{
		"variables":    env,
		"write_access": !errors.Is(err, gorm.ErrRecordNotFound),
		"owner_id":     repo.UserID,
		"is_owner":     repo.UserID == user.ID,
		"contributors": contributors,
	}

	api.NewResponse(w).
		Header("Cache-Control", "max-age=10").
		Status(http.StatusOK).
		JSON(&response)
}

type contributor struct {
	Login       string `json:"login"`
	ID          uint64 `json:"id"`
	AvatarUrl   string `json:"avatar_url"`
	WriteAccess bool   `json:"write_access"`
}

func getContributors(wg *sync.WaitGroup, token, repoFullName string, c *[]contributor, err error) {
	wg.Add(1)
	defer wg.Done()

	var res *http.Response
	res, err = gh.New(token).Get("/repos/%s/collaborators", repoFullName)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(c); err != nil {
		return
	}
}
