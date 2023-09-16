package variables

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

func handleUnlink(w http.ResponseWriter, r *http.Request) {
	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 64)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error("Invalid repository id.")
		return
	}

	db := database.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `DELETE FROM repositories WHERE user_id = ? AND id = ?;`
	result, err := db.ExecContext(ctx, q, user.ID, repoID)
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if rows == 0 {
		api.NewResponse(w).Status(http.StatusForbidden).Error("You are not the repository's owner.")
		return
	}

	// Respond with success
	api.NewResponse(w).Status(http.StatusNoContent).Done()
}
