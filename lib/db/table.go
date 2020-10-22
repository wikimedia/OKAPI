package db

import "strings"

// Table struct to represent database table
type Table struct {
	Name        string
	PrimaryKey  []string
	Columns     []Column
	Indexes     []Index
	ForeignKeys []ForeignKey
	Partition   *Partition
}

// Create function to create table sql
func (table *Table) Create() string {
	sql := "create table " + table.Name + " ("

	for _, column := range table.Columns {
		sql += column.Define() + ","
	}

	if len(table.PrimaryKey) > 0 {
		sql += "PRIMARY KEY (" + strings.Join(table.PrimaryKey, ",") + ")"
	}

	sql = strings.Trim(sql, `,`)

	if table.Partition != nil {
		sql += ") PARTITION by " + string(table.Partition.By) + "(" + table.Partition.Field + ");"
	} else {
		sql += ");"
	}

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
