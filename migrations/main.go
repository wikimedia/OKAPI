package main

import (
	"os"

	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	"gopkg.in/gookit/color.v1"
	"okapi/lib/db"
	"okapi/lib/env"
)

func main() {
	env.Context.Parse(".env")
	db := db.Client()
	defer db.Close()

	err := migrations.Run(db, "./", os.Args)

	if err != nil {
		color.Error.Println(err)
	}
}
