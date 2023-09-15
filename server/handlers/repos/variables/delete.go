package variables

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"gorm.io/gorm"
)

func Delete(w http.ResponseWriter, r *http.Request) {
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
	log.Println(user)

	variableID := r.URL.Query().Get("id")
	if variableID == "" {
		api.NewResponse(w).Status(http.StatusBadRequest).Error("invalid param: id")
		return
	}

	// check for user's write access
	// err = db.Where("repository_id = ? AND user_id = ?", repoID, user.ID).
	// 	First(&database.Permission{}).
	// 	Error
	switch err {
	case gorm.ErrRecordNotFound:
		api.NewResponse(w).
			Status(http.StatusForbidden).
			Error("you don't have write access, please contact repository owner.")
		return

	case nil:
		break

	default:
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// delete
	err = deleteVariable(variableID, repoID)
	switch err {
	case gorm.ErrRecordNotFound:
		api.NewResponse(w).Status(http.StatusNotFound).Error("variable not found.")
		return

	case nil:
		api.NewResponse(w).Status(http.StatusNoContent).Done()

	default:
		api.NewResponse(w).ServerError(err.Error())
		return
	}
}

func deleteVariable(varID string, repoID uint64) error {
	return nil
	/*
		return db.Transaction(func(tx *gorm.DB) error {
			result := tx.Where("id = ?", varID).Delete(&database.Variable{})
			if err := result.Error; err != nil {
				return err
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}

			var repo database.Repository
			err := tx.Where("id = ?", repoID).First(&repo).Error
			if err != nil {
				return err
			}

			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}

			return nil
		})
	*/
}
