package db

// Column struct to represent column in database
type Column struct {
	TableName string
	Name      string
	Type      string
}

// Add function to add new column
func (column *Column) Add() string {
	return "ALTER TABLE " + column.TableName + " ADD COLUMN " + column.Define() + ";"
}

// Set set not null to the column
func (column *Column) Set(value string) string {
	return "ALTER TABLE " + column.TableName + " ALTER COLUMN " + column.Name + " SET " + value + ";"
}

// Define function to get column definition
func (column *Column) Define() string {
	return column.Name + " " + column.Type
}

// Drop function to drop column in db
func (column *Column) Drop() string {
	return "ALTER TABLE " + column.TableName + " DROP COLUMN " + column.Name + ";"
}
