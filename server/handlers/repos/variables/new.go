package variables

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// request body: { key: string, value: string }
func NewVariable(w http.ResponseWriter, r *http.Request) {
	// VALIDATE REQUEST
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	rCtx, ok := r.Context().Value("repoCtx").(*RepoContext)
	if !ok {
		api.NewResponse(w).ServerError(errors.New("invalid repo context"))
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
	variable.RepositoryID = rCtx.RepoID
	if err := variable.GenID(); err != nil {
		api.NewResponse(w).ServerError(err)
	}
	if err := variable.EncryptValue(); err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	variable.CreatedAt = database.TimeNow()
	variable.UpdatedAt = database.TimeNow()

	err := db.newVariable(&variable)
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
}
