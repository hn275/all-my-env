package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hn275/envhub/server/db"
	"gorm.io/gorm"
)

var users []db.User = []db.User{
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

var repos []db.Repository = []db.Repository{
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
		Url:       "github.com/bar/test",
		UserID:    2,
	},
	{
		ID:        3,
		CreatedAt: db.TimeNow(),
		FullName:  "baz/testbaz",
		Url:       "github.com/baz/test",
		UserID:    3,
	},
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing args")
	}

	action := os.Args[1]

	switch action {
	case "seed":
		seed()
	default:
		fmt.Println("usage: [seed]")
	}
}

func seed() {
	d := db.New()

	defer func(*gorm.DB) {
		it, ok := recover().(error)
		if !ok {
			return
		}
		fmt.Fprint(os.Stderr, it.Error())
		d.Delete(users)
		d.Delete(repos)
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
}
