package variables

import (
	"errors"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type models interface {
	getRepoAccess(uint64, uint64) error
}

type variableDB struct {
	*gorm.DB
}

var (
	db *variableDB

	errRepoIDNotFound = errors.New("repository id not found.")
)

// `getRepoAccess` is used to check for write access. requires the `repo.ID` to
// be set, and will marshal the `repo.FullName` field
//
// returned error:
//  1. `repo.ID` is not set (equals 0): `errRepoIDNotFound`
//  2. if the user with `userID` does not have write access: `gorm.ErrRecordNotFound`
func (db *variableDB) getRepoAccess(repo *database.Repository, userID uint64) error {
	if repo.ID == 0 {
		return errRepoIDNotFound
	}

	return db.Table(database.TableRepos).
		Select("repositories.full_name").
		Where("permissions.repository_id = ? AND users.id = ?", repo.ID, userID).
		InnerJoins("INNER JOIN permissions ON permissions.repository_id = repositories.id").
		InnerJoins("INNER JOIN users ON permissions.user_id = users.id").
		First(&repo).Error
}
