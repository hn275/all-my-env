package variables

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

// request body: { key: string, value: string }
func NewVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// marshal and validate request body
	var variable database.Variable
	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	v := validator.New()
	if err := v.Struct(&variable); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("Missing one or more required fields: key, value.")
		return
	}

	// SERIALIZE VARIABLE
	variable.RepositoryID = uint32(repoID)

	if err := variable.GenID(); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if err := variable.EncryptValue(); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	variable.CreatedAt = time.Now().UTC()

	// CHECK FOR WRITE ACCESS
	wa, err := db.hasWriteAccess(user.ID, uint32(repoID))
	if !wa {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("you do not have write access, please contact the repository owner.")
		return
	}

	// WRITE TO DB
	// time
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `
		INSERT INTO variables (id,created_at,variable_key,variable_value,repository_id) 
			VALUES (:id,:created_at,:variable_key,:variable_value,:repository_id);
	`

	_, err = db.NamedExecContext(ctx, q, &variable)
	if err == nil {
		api.NewResponse(w).Status(http.StatusCreated).JSON(&variable)
		return
	}

	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	switch sqlErr.Number {

	case database.ErrDuplicateEntry:
		api.NewResponse(w).
			Status(http.StatusConflict).
			Error("Variable `%s` exists", variable.Key)
		return

	default:
		os.Stderr.Write([]byte(sqlErr.Error()))
		api.NewResponse(w).
			Status(http.StatusInternalServerError).
			Error("Failed to insert to database, try again later.")
		return
	}
}
