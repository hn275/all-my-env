package variables

import (
	"errors"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type Model struct {
	*gorm.DB
}

var db *Model

func init() {
	db = &Model{database.New()}
}

func (db *Model) getVariables(v *[]database.Variable, repoID uint64) error {
	return db.Table(database.TableVariables).
		Where("repository_id = ?", repoID).
		Find(&v).Error
}

func (db *Model) getWriteAccess(userID, repoID uint64, perm *database.Permission) error {
	err := db.
		Where("user_id = ? AND repository_id = ?", userID, repoID).
		First(&perm).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (db *Model) newVariable(v *database.Variable) error {
	return db.Create(v).Error
}
