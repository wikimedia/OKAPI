package main

import (
	"log"
	"okapi-data-service/lib/env"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

const directory = "./migrations"

func main() {
	if err := env.Init(); err != nil {
		log.Panic(err)
	}

	db := pg.Connect(&pg.Options{
		MaxRetries:      5,
		MinRetryBackoff: 2 * time.Second,
		MaxRetryBackoff: 10 * time.Second,
		Addr:            env.DBAddr,
		User:            env.DBUser,
		Database:        env.DBName,
		Password:        env.DBPassword,
	})

	defer db.Close()

	if err := migrations.Run(db, directory, os.Args); err != nil {
		log.Panic(err)
	}
}
