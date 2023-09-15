package variables

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

// method: PUT
// jsson body: {key: string, value: string}
func handleEdit(w http.ResponseWriter, r *http.Request) {
	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	variableID := chi.URLParam(r, "variableID")
	if variableID == "" {
		api.NewResponse(w).Status(http.StatusBadRequest).Error("Variable ID not found.")
		return
	}

	var variable database.Variable
	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// check access
	wa, err := db.hasWriteAccess(user.ID, uint32(repoID))
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if !wa {
		api.
			NewResponse(w).
			Status(http.StatusForbidden).
			Error("you don't have write access to this repository, please contact the owner.")
		return
	}

	// serialize
	variable.RepositoryID = uint32(repoID)
	variable.ID = variableID
	variable.UpdatedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}

	err = variable.EncryptValue()
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// updates and writes back to db
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `
	UPDATE variables SET variable_key = ?, variable_value = ?, updated_at = ?
	WHERE id = ? AND repository_id = ?;
	`

	result, err := db.ExecContext(
		ctx,
		q,
		variable.Key,
		variable.Value,
		variable.UpdatedAt,
		variable.ID,
		variable.RepositoryID,
	)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	row, err := result.RowsAffected()
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if row == 0 {
		api.NewResponse(w).Status(http.StatusNotFound).Error("Variable not found.")
		return
	}

	api.NewResponse(w).Status(http.StatusOK).JSON(map[string]time.Time{"updated_at": variable.UpdatedAt.Time})
}
