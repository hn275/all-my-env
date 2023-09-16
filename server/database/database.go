package database

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/hn275/envhub/server/lib"
	"github.com/jmoiron/sqlx"
)

var (
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
