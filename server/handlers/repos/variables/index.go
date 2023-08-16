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
	type Repo struct {
		database.Repository
		Login string
	}
	var repo Repo
	repo.ID, err = strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Invalid repository id: %s", err.Error())
		return
	}
	sel := []string{
		"users.login",
		"repositories.full_name",
		"repositories.url",
		"repositories.variable_count",
	}
	db := database.New()
	err = db.Table(database.TableRepos).
		Select(sel).
		InnerJoins("JOIN users ON users.id = repositories.user_id").
		Where("repositories.id = ? AND users.id = ?", repo.ID, user.ID).
		First(&repo).Error
	switch err {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).Status(http.StatusNotFound).Error("Repository not found")
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
		userLogin: repo.Login,
		userTok:   user.Token,
		repoURL:   repo.FullName,
	}
	go c.getRepoAccess()

	var env []database.Variable
	err = db.Table(database.TableVariables).
		Where("repository_id = ?", repo.ID).
		Find(&env).Error
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

	api.NewResponse(w).Status(http.StatusOK).JSON(&env)
}
