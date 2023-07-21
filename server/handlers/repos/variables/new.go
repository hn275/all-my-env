package variables

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
)

func (d *variableHandler) NewVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	s := chi.URLParam(r, "id")
	repoID, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		api.NewResponse(w).
			Status(http.StatusBadRequest).
			Error("invalid id url param")
		return
	}

	api.NewResponse(w).Status(http.StatusOK).Text(string(rune(repoID)))
}
