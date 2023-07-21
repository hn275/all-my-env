package variables

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/db"
)

func (d *variableHandler) NewVariable(w http.ResponseWriter, r *http.Request) {
	// VALIDATE REQUEST
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	s := chi.URLParam(r, "id")
	repoID, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	var body EnvVariable
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	// SERIALIZE VARIABLE
	envVar, err := body.Cipher(uint32(repoID))
	if err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}
	envVar.CreatedAt = db.TimeNow()
	envVar.UpdatedAt = db.TimeNow()

	// SAVE TO DB
	// TODO: handle duplicate env var key
	if err := d.Create(envVar).Error; err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	api.NewResponse(w).Status(http.StatusCreated).Done()
}
