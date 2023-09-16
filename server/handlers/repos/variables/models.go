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
