package repos

import (
	"log"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type repoModels interface {
	newRepo(repoBuf *database.Repository) error
	findRepo(userID uint64, ids []uint64) ([]database.Repository, error)
	deleteRepo(repoID uint64, userID uint64) error
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

func (db *repoDatabase) deleteRepo(repoID uint64, userID uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete the repository
		log.Println("Deleting repository with repoID:", repoID)
		result := tx.Where("id = ? AND user_id = ?", repoID, userID).Delete(&database.Repository{})
		if err := result.Error; err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		// Delete the associated variables
		result = tx.Where("repository_id = ?", repoID).Delete(&database.Variable{})
		if err := result.Error; err != nil {
			return err
		}

		return nil
	})
}
