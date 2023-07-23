package variables

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"gorm.io/gorm"
)

type Repository struct {
	db.Repository
	Variable []db.Variable `json:"variables"`
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

	repoID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Invalid repository id: %s", err.Error())
		return
	}

	// QUERY DB FOR REPO INFO
	var repo Repository
	err = h.Table(db.TableRepos).Where("id = ?", repoID).First(&repo).Error

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
	}

	// CHECK FOR USER ACCESS
	c := make(chan permission)
	defer close(c)
	go getRepoAccess(c, repo.FullName, user)

	// QUERY DB FOR ENV VARIABLES
	err = h.Model(&[]db.Variable{}).
		Where("repository_id = ?", repoID).
		Find(&repo.Variable).Error // TODO: add pagination
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		api.NewResponse(w).ServerError(err)
		return
	}

	// decrypt values
	for i := range repo.Variable {
		err := repo.Variable[i].DecryptValue()
		if err != nil {
			api.NewResponse(w).ServerError(err)
			return
		}
	}

	for {
		select {
		case a := <-c:
			if a.err != nil {
				api.NewResponse(w).
					Status(http.StatusBadGateway).
					Error(a.err.Error())
				return
			}

			if !a.allowed {
				api.NewResponse(w).
					Status(http.StatusForbidden).
					Done()
			}

			// debuggo
			j, err := json.MarshalIndent(repo, "", "  ")
			if err != nil {
				panic(err)
			}
			log.Fatal(string(j))
			// debuggo

			api.NewResponse(w).Status(http.StatusOK).JSON(repo)
			return
		default:
			continue
		}
	}
}
