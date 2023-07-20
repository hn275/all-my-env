package db

import (
	"fmt"
	"log"
	"time"

	"github.com/hn275/envhub/server/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string

	db  *gorm.DB
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

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// auto migrate
	autoMigrate(db, &User{})
	autoMigrate(db, &Repository{})
	autoMigrate(db, &Variable{})
	fmt.Println("Automigrate done")
}

func New() *gorm.DB {
	return db
}

func autoMigrate(db *gorm.DB, d interface{}) {
	if err := db.AutoMigrate(&d); err != nil {
		log.Fatal(err)
	}
}

// returns UTC time, ie: 2023-07-06T07:25:26Z
func TimeNow() TimeStamp {
	return time.Now().UTC()
}
