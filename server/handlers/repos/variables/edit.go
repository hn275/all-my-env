package variables

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

func Edit(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	var variable database.Variable
	if err := json.NewDecoder(r.Body).Decode(&variable); err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// serialize
	var serializeErr error
	wg := sync.WaitGroup{}
	go serializeVariable(&wg, &variable, serializeErr)

	// check access
	writeAccess, err := db.hasWriteAccess(user.ID, repoID)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if !writeAccess {
		api.
			NewResponse(w).
			Status(http.StatusForbidden).
			Error("you don't have write access to this repository, please contact the owner.")
		return
	}

	// updates and writes back to db
	wg.Wait()

	u := map[string]any{
		"key":        variable.Key,
		"value":      variable.Value,
		"updated_at": database.TimeNow(),
	}
	result := db.Model(&variable).
		Where("id = ? AND repository_id = ?", variable.ID, repoID).
		Updates(u)

	if result.Error != nil {
		api.NewResponse(w).ServerError(result.Error.Error())
		return
	}

	if result.RowsAffected == 0 {
		api.NewResponse(w).Status(http.StatusNotFound).Error("variable not found.")
		return
	}

	delete(u, "value")
	delete(u, "key")
	api.NewResponse(w).Status(http.StatusOK).JSON(u)
}

func serializeVariable(wg *sync.WaitGroup, v *database.Variable, err error) {
	wg.Add(1)
	defer wg.Done()
	err = v.EncryptValue()
	if err != nil {
		return
	}
}
