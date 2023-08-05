package variables

import (
	"errors"

	"github.com/hn275/envhub/server/database"
	"gorm.io/gorm"
)

type variableHandler struct {
	*gorm.DB
}

var (
	Handlers      *variableHandler
	errBadGateWay = errors.New("GitHub responded an with error")
)

func init() {
	Handlers = &variableHandler{database.New()}
	db = &variableDB{database.New()}
}
