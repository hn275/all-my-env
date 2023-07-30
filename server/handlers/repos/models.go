package repos

import (
	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

var (
	db repoModels
)

type repoModels interface {
	newRepo(*database.Repository) error
}

type repoDatabase struct{ *gorm.DB }

func init() {
	db = &repoDatabase{database.New()}
}

func (repoDB *repoDatabase) newRepo(r *database.Repository) error {
	return repoDB.Create(r).Error
}
