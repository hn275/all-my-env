package database

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hn275/envhub/server/lib"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

const (
	TableRepos       = ""
	TablePermissions = ""
	TableUsers       = ""
	TableVariables   = ""
)

var (
	db  *gorm.DB
	dbx *sqlx.DB
	err error
)

func init() {
	// sqlx
	mysqlDsn := lib.Getenv("MYSQL_DSN")
	dbx = sqlx.MustConnect("mysql", mysqlDsn)
}

func New() *sqlx.DB {
	return dbx
}

func autoMigrate(db *gorm.DB, d interface{}) {
	if err := db.AutoMigrate(&d); err != nil {
		log.Fatal(err)
	}
}

type TimeStamp = time.Time

// returns UTC time, ie: 2023-07-06T07:25:26Z
func TimeNow() TimeStamp {
	return time.Now().UTC()
}
