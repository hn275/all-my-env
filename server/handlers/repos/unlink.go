package repos

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
)

// returns 204 on success, no body
// { "message": "err" } otherwise
func Unlink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

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
	result := db.Where("id = ? AND user_id = ?", repoID, user.ID).Delete(&database.Repository{})
	if err := result.Error; err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if result.RowsAffected == 0 {
		api.NewResponse(w).Error("Not repository owner")
		return
	}

	// Respond with success
	api.NewResponse(w).Status(http.StatusNoContent).Done()
}
