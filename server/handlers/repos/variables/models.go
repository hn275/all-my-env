package variables

import (
	"errors"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type model struct {
	*gorm.DB
}

var db *model

func init() {
	db = &model{database.NewGorm()}
}

func (db *model) getVariables(v *[]database.Variable, repoID uint64) error {
	return db.Table(database.TableVariables).
		Where("repository_id = ?", repoID).
		Find(&v).Error
}

// deprecated: no need to pass in the additional `perm` param, use `hasWriteAccess` instead
func (db *model) getWriteAccess(userID, repoID uint64, perm *database.Permission) error {
	err := db.
		Where("user_id = ? AND repository_id = ?", userID, repoID).
		First(&perm).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (db *model) hasWriteAccess(userID, repoID uint32) (bool, error) {
	err := db.
		Where("user_id = ? AND repository_id = ?", userID, repoID).
		First(&database.Permission{}).
		Error

	switch err {
	case nil:
		return true, nil

	case gorm.ErrRecordNotFound:
		return false, nil

	default:
		return false, err
	}
}
