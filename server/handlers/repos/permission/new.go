package permission

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hn275/envhub/server/api"
	"github.com/hn275/envhub/server/database"
	"github.com/jmoiron/sqlx"
)

// request json contains all the user ids that should have write access:
// { "userIDs": []uint64 }
func NewPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.NewResponse(w).Status(http.StatusMethodNotAllowed).Done()
		return
	}

	user, err := api.NewContext(r).User()
	if err != nil {
		api.NewResponse(w).Status(http.StatusForbidden).Error(err.Error())
		return
	}

	var newPerm struct {
		UserIDs []uint32 `json:"userIDs"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPerm); err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
		return
	}

	// check to for repository owner.
	repoID, err := strconv.ParseUint(chi.URLParam(r, "repoID"), 10, 32)
	if err != nil {
		api.NewResponse(w).Status(http.StatusBadRequest).Error(err.Error())
	}

	db := database.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var repoOwner uint32
	q := `SELECT user_id FROM repositories WHERE id = ?;`
	err = db.QueryRowxContext(ctx, q, repoID).Scan(&repoOwner)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.NewResponse(w).Status(http.StatusForbidden).Error("You are not repository owner")
			return
		}
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// get existing users with write access
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var waUsers []uint32
	q = `SELECT user_id FROM permissions WHERE repository_id = ? AND NOT user_id = ?;`
	row, err := db.QueryxContext(ctx, q, repoID, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		api.NewResponse(w).ServerError(err.Error())
		return
	}
	for row.Next() {
		var i uint32
		if err := row.Scan(&i); err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
		waUsers = append(waUsers, i)
	}

	// get user ids without write access
	revokeWa := make([]uint32, 0, len(waUsers))
	for _, userID := range waUsers {
		if !contains(newPerm.UserIDs, userID) {
			revokeWa = append(revokeWa, userID)
		}
	}

	// update database
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})

	if err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	// delete revoked-access entries
	revokeWa = append(revokeWa, 1234)
	if len(revokeWa) != 0 {
		q = `DELETE FROM permissions WHERE repository_id = ? AND user_id IN (?);`
		deleteQuery, args, err := sqlx.In(q, repoID, revokeWa)
		if err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}

		_, err = tx.ExecContext(ctx, deleteQuery, args...)
		if err != nil {
			api.NewResponse(w).ServerError(err.Error())
			return
		}
	}

	// insert to db
	// NOTE: since all id's are uint, safe for string interpolation
	param := make([]string, len(newPerm.UserIDs))
	for i, userID := range newPerm.UserIDs {
		param[i] = fmt.Sprintf("(%d,%d)", repoID, userID)
	}

	q = `
	INSERT INTO permissions (repository_id, user_id) VALUES %s as new
	ON DUPLICATE KEY UPDATE repository_id = new.repository_id;
	`
	q = fmt.Sprintf(q, strings.Join(param, ","))

	if _, err := tx.ExecContext(ctx, q); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		api.NewResponse(w).ServerError(err.Error())
		return
	}

	api.NewResponse(w).Status(http.StatusOK).Done()
}

func contains(src []uint32, i uint32) bool {
	for _, j := range src {
		if i == j {
			return true
		}
	}
	return false
}
