package db

import (
	"github.com/go-pg/pg/v9"
	"okapi/lib/env"
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
