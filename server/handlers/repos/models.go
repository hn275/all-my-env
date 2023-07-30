package repos

import (
	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

var (
	database repoDB
)

type repoDB interface {
	newRepo(*db.Repository) error
}

type repoDatabase struct{ *gorm.DB }

func init() {
	database = &repoDatabase{db.New()}
}

func (repoDB *repoDatabase) newRepo(r *db.Repository) error {
	return repoDB.Create(r).Error
}
