package models

import (
	"fmt"
	pg_db "okapi/lib/db"

	"github.com/go-pg/pg/v9"
)

var db *pg.DB

// DB Get database client
func DB() *pg.DB {
	if db == nil {
		db = pg_db.Client()
	}

	return db
}

// Init function to init db on startup
func Init() error {
	if DB() == nil {
		return fmt.Errorf("Database client not created")
	}

	_, err := db.Exec("select 1")

	return err
}
