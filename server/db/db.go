package db

import (
	"fmt"
	"log"
	"os"

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
	db       *sqlx.DB
	err      error
)

func init() {
	host = getEnv("POSTGRES_HOST")
	user = getEnv("POSTGRES_USER")
	password = getEnv("POSTGRES_PASSWORD")
	dbname = getEnv("POSTGRES_DB")
	port = getEnv("POSTGRES_PORT")
	sslmode = getEnv("POSTGRES_SSLMODE")

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

func getEnv(k string) string {
	t := os.Getenv(k)
	if t == "" {
		log.Fatalf("[%s] not set", k)
	}

	return t
}

func Close() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
