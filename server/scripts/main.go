package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

var (
	d *gorm.DB
)

func init() {
	// clean db for seeding
	d = db.New()

	// clean db
	d.Raw("DELETE FROM users")
	d.Raw("DELETE FROM repositories")
	d.Raw("DELETE FROM variables")
	d.Raw("DELETE FROM permissions")
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing args")
	}

	action := os.Args[1]

	switch action {
	case "mock":
		mock()
	default:
		fmt.Println("usage: [mock]")
	}
}

func mock() {
	users := []db.User{
		{
			ID:        1,
			CreatedAt: db.TimeNow(),
			Vendor:    "github",
			UserName:  "foo",
		},
		{
			ID:        2,
			CreatedAt: db.TimeNow(),
			Vendor:    "github",
			UserName:  "bar",
		},
		{
			ID:        3,
			CreatedAt: db.TimeNow(),
			Vendor:    "github",
			UserName:  "baz",
		},
	}

	repos := []db.Repository{
		{
			ID:        1,
			CreatedAt: db.TimeNow(),
			FullName:  "foo/testfoo",
			Url:       "github.com/foo/testfoo",
			UserID:    1,
		},
		{
			ID:        2,
			CreatedAt: db.TimeNow(),
			FullName:  "bar/testbar",
			Url:       "github.com/bar/testbar",
			UserID:    2,
		},
		{
			ID:        3,
			CreatedAt: db.TimeNow(),
			FullName:  "baz/testbaz",
			Url:       "github.com/baz/testbaz",
			UserID:    3,
		},
	}

	perms := db.Permission{
		ID:           1,
		RepositoryID: 1,
		UserID:       1,
	}

	// VARIABLE
	variable := db.Variable{
		ID:           "",
		CreatedAt:    db.TimeNow(),
		UpdatedAt:    db.TimeNow(),
		Key:          "foo",
		Value:        "bar",
		Repository:   db.Repository{},
		RepositoryID: 1,
	}
	variable.GenID()
	if err := variable.EncryptValue(); err != nil {
		panic(err)
	}

	defer func(*gorm.DB) {
		it, ok := recover().(error)
		if !ok {
			return
		}
		fmt.Fprint(os.Stderr, it.Error())
		d.Delete(users)
		d.Delete(repos)
		d.Delete(&perms)
		d.Delete(&variable)
	}(d)

	for _, u := range users {
		err := d.Create(&u).Error
		if err != nil {
			panic(err)
		}
	}

	for _, r := range repos {
		err := d.Create(&r).Error
		if err != nil {
			panic(err)
		}
	}

	err := d.Create(&perms).Error
	if err != nil {
		panic(err)
	}

	err = d.Create(&variable).Error
	if err != nil {
		panic(err)
	}
}
