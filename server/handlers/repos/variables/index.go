package variables

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type Repository struct {
	Meta      db.Repository `json:",inline"`
	Variables []db.Variable `json:"variables"`
}

func (h *variableHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := jwt.GetUser(r)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Invalid repository id: %s", err.Error())
		return
	}

	// QUERY DB FOR REPO INFO
	var repo Repository
	err = h.Table(db.TableRepos).Where("id = ?", repoID).First(&repo.Meta).Error

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		api.NewResponse(w).
			Status(http.StatusNotFound).
			Error("Repository not found")
		return
	default:
		api.NewResponse(w).ServerError(err)
		return
	}

	// CHECK FOR USER ACCESS
	c := make(chan error, 1)
	wg := new(sync.WaitGroup)
	defer close(c)
	go getRepoAccess(c, wg, repo.Meta.FullName, user)

	// QUERY DB FOR ENV VARIABLES
	err = h.Model(&[]db.Variable{}).
		Where("repository_id = ?", repoID).
		Find(&repo.Variables).Error // TODO: add pagination
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err)
		return
	}

	// decrypt values
	for i := range repo.Variables {
		err := repo.Variables[i].DecryptValue()
		if err != nil {
			api.NewResponse(w).ServerError(err)
			return
		}
	}

	wg.Wait()
	switch err := <-c; err {
	case nil:
		api.NewResponse(w).Status(http.StatusOK).JSON(repo)
		return
	case errNotAContributor:
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return

	case errBadGateWay:
		api.NewResponse(w).Status(http.StatusBadGateway).Error(err.Error())
		return

	default:
		api.NewResponse(w).ServerError(err)
	}

}
