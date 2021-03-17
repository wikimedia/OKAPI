package main

import (
	"log"
	"okapi-data-service/lib/env"
	"os"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "/migrations"

func main() {
	err := env.Init()

	if err != nil {
		log.Panic(err)
	}

	db := pg.Connect(&pg.Options{
		Addr:     env.DBAddr,
		User:     env.DBUser,
		Database: env.DBName,
		Password: env.DBPassword,
	})

	err = migrations.Run(db, directory, os.Args)

	if err != nil {
		log.Panic(err)
	}
}
