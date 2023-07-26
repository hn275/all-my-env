package variables

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
	jwt "github.com/hn275/envhub/server/jsonwebtoken"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type permission struct {
	allowed bool
	err     error
}

func (d *variableHandler) NewVariable(w http.ResponseWriter, r *http.Request) {
	// VALIDATE REQUEST
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := jwt.GetUser(r)
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Done()
		return
	}

	var body db.Variable
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("invalid repository id")
		return
	}

	// CHECKS IF USER IS A COLLABORATOR
	var repo struct {
		FullName string
	}

	err = d.Table(db.TableRepos).
		Select("repositories.full_name").
		Where("permissions.repository_id = ? AND users.id = ?", repoID, user.ID).
		InnerJoins("INNER JOIN permissions ON permissions.repository_id = repositories.id").
		InnerJoins("INNER JOIN users ON permissions.user_id = users.id").
		First(&repo).Error

	switch err {
	case nil:
		break

	case gorm.ErrRecordNotFound:
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Write-access not granted. Please contact repo owner.")
		return

	default:
		api.NewResponse(w).ServerError(err)
		return
	}

	ghChan := make(chan permission)
	go getRepoAccess(ghChan, repo.FullName, user)

	// SERIALIZE VARIABLE
	body.RepositoryID = uint32(repoID)
	body.GenID()
	if err := body.EncryptValue(); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	body.CreatedAt = db.TimeNow()
	body.UpdatedAt = db.TimeNow()

	// WRITE TO DB
	for {
		select {
		case access := <-ghChan:
			defer close(ghChan)
			if access.err != nil {
				api.NewResponse(w).
					Status(http.StatusBadGateway).
					Error(access.err.Error())
				return
			}

			if !access.allowed {
				api.NewResponse(w).Status(http.StatusForbidden).Done()
				return
			}

			err := d.Create(&body).Error
			if err == nil {
				api.NewResponse(w).Status(http.StatusCreated).Done()
				return
			}

			pgErr, ok := err.(*pgconn.PgError)
			if !ok {
				api.NewResponse(w).ServerError(err)
				return
			}

			if pgErr.Code == pgerrcode.UniqueViolation {
				api.NewResponse(w).
					Status(http.StatusConflict).
					Error(pgErr.Error())
				return
			}

			api.NewResponse(w).ServerError(pgErr)
			return

		default:
			continue
		}
	}
}
