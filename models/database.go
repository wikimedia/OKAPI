package models

import (
	"github.com/go-pg/pg/v9"
	"okapi/lib/db"
)

// Database database instance
var database *pg.DB

// DB Get database client
func DB() *pg.DB {
	if database == nil {
		database = db.Client()
	}

	return database
}
