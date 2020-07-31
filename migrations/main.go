package main

import (
	"os"

	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	"okapi/lib/db"
	"okapi/lib/env"
)

func main() {
	env.Context.Parse(".env")
	db := db.Client()
	defer db.Close()
	migrations.Run(db, "./", os.Args)
}
