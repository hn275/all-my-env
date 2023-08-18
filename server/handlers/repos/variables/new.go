package variables

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// request body: { key: string, value: string }
func NewVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	var variable database.Variable
	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	// SERIALIZE VARIABLE
	variable.RepositoryID = repoID
	// gen id
	if err := variable.GenID(); err != nil {
		api.NewResponse(w).ServerError(err.Error())
	}
	// cipher value
	if err := variable.EncryptValue(); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	// time
	variable.CreatedAt = database.TimeNow()
	variable.UpdatedAt = database.TimeNow()

	err = db.newVariable(&variable)
	if err == nil {
		api.NewResponse(w).Status(http.StatusCreated).JSON(&variable)
		return
	}

	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if pgErr.Code == pgerrcode.UniqueViolation {
		api.NewResponse(w).
			Status(http.StatusConflict).
			Error(pgErr.Error())
		return
	}

	api.NewResponse(w).ServerError(pgErr.Error())
}
