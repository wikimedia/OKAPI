package db

import (
	"okapi/lib/env"

	"github.com/go-pg/pg/v10"
)

// Client returns new db connection (don't forget to close it with defer)
func Client() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     env.Context.DBAddr,
		User:     env.Context.DBUser,
		Database: env.Context.DBName,
		Password: env.Context.DBPassword,
	})
}
