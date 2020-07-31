package db

import "strings"

// Table struct to represent database table
type Table struct {
	Name        string
	Columns     []Column
	Indexes     []Index
	ForeignKeys []ForeignKey
}

// Create function to create table sql
func (table *Table) Create() string {
	sql := "create table " + table.Name + " ("

	for _, column := range table.Columns {
		sql += column.Define() + ","
	}

	sql = strings.Trim(sql, `,`) + ");"

	for _, index := range table.Indexes {
		sql += index.Create()
	}

	for _, fk := range table.ForeignKeys {
		sql += fk.Create()
	}

	return sql
}

// Drop function to create drop table sql
func (table *Table) Drop() string {
	sql := ""

	for _, fk := range table.ForeignKeys {
		sql += fk.Drop()
	}

	for _, index := range table.Indexes {
		sql += index.Drop()
	}

	return sql + " DROP TABLE " + table.Name + ";"
}
