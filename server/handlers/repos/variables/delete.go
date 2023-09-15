package variables

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
)

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
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
		api.NewResponse(w).Status(http.StatusBadRequest).Error("Missing variable ID not found.")
		return
	}

	// check for user's write access
	wa, err := db.hasWriteAccess(user.ID, uint32(repoID))
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if !wa {
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("You do not have write access to this repository, please contact your repository owner.")
		return
	}

	// delete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `DELETE FROM variables WHERE id = ? AND repository_id = ?`
	_, err = db.ExecContext(ctx, q, variableID, repoID)

	switch err {
	case nil:
		api.NewResponse(w).Status(http.StatusNoContent).Done()
		return

	case sql.ErrNoRows:
		api.NewResponse(w).Status(http.StatusNotFound).Error("Variable not found.")
		return

	default:
		api.NewResponse(w).ServerError(err.Error())
	}

}
