package main

import (
	"os"

	"okapi/lib/env"

	"github.com/go-pg/pg/v9"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	"gopkg.in/gookit/color.v1"
)

func main() {
	env.Context.Parse(".env")
	db := pg.Connect(&pg.Options{
		Addr:     env.Context.DBAddr,
		User:     env.Context.DBUser,
		Database: env.Context.DBName,
		Password: env.Context.DBPassword,
	})
	defer db.Close()

	err := migrations.Run(db, "./", os.Args)

	if err != nil {
		color.Error.Println(err)
	}
}
