package variables

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hn275/envhub/server/database"
	"github.com/jmoiron/sqlx"
)

type model struct {
	*sqlx.DB
}

var db *model

func init() {
	db = &model{database.New()}
}

func (db *model) getVariables(v *[]database.Variable, repoID uint64) error {
	// return db.Table(database.TableVariables).
	// 	Where("repository_id = ?", repoID).
	// 	Find(&v).Error
	return nil
}

// deprecated: no need to pass in the additional `perm` param, use `hasWriteAccess` instead
func (db *model) getWriteAccess(userID, repoID uint64, perm *database.Permission) error {
	return nil
	// err := db.
	// 	Where("user_id = ? AND repository_id = ?", userID, repoID).
	// 	First(&perm).
	// 	Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil
	// }
	// return err
}

func (db *model) hasWriteAccess(userID, repoID uint32) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q := `SELECT (id) FROM permissions WHERE (user_id = ? AND repository_id = ?);`
	_, err := db.QueryxContext(ctx, q, userID, repoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
