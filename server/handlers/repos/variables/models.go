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

type RepoInfo struct {
	*database.Repository
	UserLogin string `gorm:"column:login"`
}

func (db *Model) getRepoInfo(r *RepoInfo) error {
	if r.UserID == 0 {
		return errors.New("user id not found.")
	}
	if r.ID == 0 {
		return errors.New("repository id not found.")
	}
	sel := []string{
		"users.login",
		"repositories.full_name",
		"repositories.url",
		"repositories.variable_count",
	}
	return db.Table(database.TableRepos).
		Select(sel).
		InnerJoins("JOIN users ON users.id = repositories.user_id").
		Where("repositories.id = ? AND users.id = ?", r.ID, r.UserID).
		First(r).Error
}

func (db *Model) getVariables(v []database.Variable, repoID uint64) error {
	return db.Table(database.TableVariables).
		Where("repository_id = ?", repoID).
		Find(&v).Error
}
