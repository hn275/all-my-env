package repos

import (
	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type repoModels interface {
	newRepo(repoBuf *database.Repository) error
	findRepo(userID uint64, ids []uint64) ([]database.Repository, error)
	deleteRepo(id uint64) error
}

type repoDatabase struct{ *gorm.DB }

var db repoModels

func init() {
	db = &repoDatabase{database.New()}
}

func (repoDB *repoDatabase) newRepo(r *database.Repository) error {
	return repoDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(r).Error; err != nil {
			return err
		}

		// enable write permission to for owner
		perm := database.Permission{
			RepositoryID: r.ID,
			UserID:       r.UserID,
		}
		if err := tx.Create(&perm).Error; err != nil {
			return err
		}
		return nil
	})
}

func (db *repoDatabase) findRepo(userID uint64, ids []uint64) ([]database.Repository, error) {
	var repos []database.Repository
	err := db.Table(database.TableRepos).
		Select([]string{"id", "variable_count"}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&repos).Error
	return repos, err
}

func (db *repoDatabase) deleteRepo(id uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&database.Repository{ID: id}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&database.Permission{RepositoryID: id}).Error; err != nil {
			return err
		}
		return nil
	})
}
