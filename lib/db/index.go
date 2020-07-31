package db

import "strings"

// Index struct to represent sql index
type Index struct {
	TableName string
	Columns   []string
}

// Create generate sql to create index
func (index Index) Create() string {
	return "CREATE INDEX " + index.TableName + "_" + strings.Join(index.Columns, "_") +
		" ON " + index.TableName + " (" + strings.Join(index.Columns, ", ") + ");"
}

// Drop generate sql to drop index
func (index Index) Drop() string {
	return "DROP INDEX " + index.TableName + "_" + strings.Join(index.Columns, "_") + ";"
}
