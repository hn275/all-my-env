package variables

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
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
	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("failed to parse repository id: %s", err.Error())
		return
	}

	repo := database.Repository{ID: repoID}
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

	// check for user access
	c := contributor{
		access:    false,
		err:       nil,
		mut:       sync.Mutex{},
		wg:        sync.WaitGroup{},
		userLogin: user.Login,
		userTok:   user.Token,
		repoURL:   repo.FullName,
	}
	go c.getRepoAccess()

	env := make([]database.Variable, repo.VariableCount)

	err = db.getVariables(&env, repo.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	for i := range env {
		err = env[i].DecryptValue()
		if err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
	}

	c.wg.Wait()
	if c.err != nil {
		api.NewResponse(w).Status(http.StatusBadGateway).Error(c.err.Error())
		return
	}

	if !c.access {
		api.NewResponse(w).Status(http.StatusForbidden).Error("not a contributor.")
		return
	}

	// query db for write access
	var p database.Permission
	err = db.getWriteAccess(repo.UserID, repo.ID, &p)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	response := map[string]any{
		"variables":    env,
		"write_access": !errors.Is(err, gorm.ErrRecordNotFound),
	}

	api.NewResponse(w).
			Header("Cache-Control", "max-age=10").
		Status(http.StatusOK).
		JSON(&response)
}
