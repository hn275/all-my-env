package variables

import (
	"errors"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type models interface {
	getRepoAccess(*database.Repository, uint64) error
	getRepoByID(uint64, *Repository) error
	getVariables(*Repository) error
}

type variableDB struct {
	*gorm.DB
}

var (
	db models

	errRepoIDNotFound = errors.New("repository id not found.")
)

func init() {
	db = &variableDB{database.New()}
}

// `getVariables` is used to fetch all variables belongs to a repo.
//
// returned errors:
//  1. `repo.ID` is not set (equals 0): `errRepoIDNotFound`
//  2. no variables found: `gorm.ErrRecordNotFound`
func (db *variableDB) getVariables(repo *Repository) error {
	if repo.ID == 0 {
		return errRepoIDNotFound
	}

	return db.Model(&[]database.Variable{}).
		Where("repository_id = ?", repo.ID).
		Find(&repo.Variables).Error
}

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

// `getRepoByID` fetches the first (and only) repo that has the `repoID` and
// marshals it into the `repo` struct
//
// returned error:
//  1. `repo.ID` is not set (equals 0): `errRepoIDNotFound`
//  2. if the repo is not found: `gorm.ErrRecordNotFound`
func (db *variableDB) getRepoByID(repoID uint64, repo *Repository) error {
	if repo.ID == 0 {
		return errRepoIDNotFound
	}

	return db.Table(database.TableRepos).
		Where("id = ?", repoID).
		First(repo).Error
}
