package pg

import (
	"errors"
	"okapi-data-service/lib/env"

	"github.com/go-pg/pg/v10"
)

// ErrDuplicateConn db connection already created
var ErrDuplicateConn = errors.New("db connection already exists")

var conn *pg.DB

// Conn get pg db connection
func Conn() *pg.DB {
	return conn
}

// Init create database connection
func Init() error {
	if conn != nil {
		return ErrDuplicateConn
	}

	conn = pg.Connect(&pg.Options{
		Addr:     env.DBAddr,
		User:     env.DBUser,
		Database: env.DBName,
		Password: env.DBPassword,
	})

	return nil
}
