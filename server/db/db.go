package db

import (
	"fmt"
	"log"

	"github.com/hn275/envhub/server/lib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string

	db  *sqlx.DB
	err error
)

func init() {
	host = lib.Getenv("POSTGRES_HOST")
	user = lib.Getenv("POSTGRES_USER")
	password = lib.Getenv("POSTGRES_PASSWORD")
	dbname = lib.Getenv("POSTGRES_DB")
	port = lib.Getenv("POSTGRES_PORT")
	sslmode = lib.Getenv("POSTGRES_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func New() *sqlx.DB {
	return db
}

// only call in main.go
func Close() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
