package repos

import (
	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type repoModels interface {
	newRepo(repoBuf *database.Repository) error
	findRepo(userID uint64, ids []uint64) ([]uint64, error)
}

type repoDatabase struct{ *gorm.DB }

var db repoModels

func init() {
	db = &repoDatabase{database.New()}
}

func (repoDB *repoDatabase) newRepo(r *database.Repository) error {
	return repoDB.Create(r).Error
}

func (db *repoDatabase) findRepo(userID uint64, ids []uint64) ([]uint64, error) {
	var repoIDs []uint64
	err := db.Table(database.TableRepos).
		Select("id").
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(repoIDs).Error
	return repoIDs, err
}
