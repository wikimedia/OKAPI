package models

import (
	pg_db "okapi/lib/db"

	"github.com/go-pg/pg/v10"
)

var db *pg.DB

// DB Get database client
func DB() *pg.DB {
	return db
}

// Init function to init db on startup
func Init() error {
	db = pg_db.Client()

	_, err := db.Exec("select 1")

	return err
}

// Close function to close db connection
func Close() {
	db.Close()
}
