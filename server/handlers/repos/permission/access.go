package permission

import (
	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = database.New()
}
