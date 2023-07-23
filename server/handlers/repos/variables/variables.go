package variables

import (
	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

type variableHandler struct {
	*gorm.DB
}

var (
	Handlers *variableHandler
)

func init() {
	Handlers = &variableHandler{db.New()}
}
