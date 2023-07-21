package variables

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
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

	var newV EnvVariable
	if err := json.NewDecoder(r.Body).Decode(&newV); err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error(err.Error())
		return
	}

	// SAVE TO DB
	v, err := newV.Serialize(uint32(repoID))
	if err := d.Create(v).Error; err != nil {
		api.NewResponse(w).ServerError(err)
		return
	}

	api.NewResponse(w).Status(http.StatusCreated).Done()
}
