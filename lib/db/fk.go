package db

// Delete postgresql schema delete function
type Delete string

// Cascade, Restrict types for OnDelete values
const (
	Cascade  Delete = "CASCADE"
	Restrict Delete = "RESTRICT"
	NoAction Delete = "NO ACTION"
)

// ForeignKey constraint struct
type ForeignKey struct {
	ParentTable string
	TableName   string
	Name        string
	Property    string
	References  string
	OnDelete    Delete
}

// Create generate SQL for creating constraint
func (foreignKey *ForeignKey) Create() string {
	return "ALTER TABLE " + foreignKey.TableName +
		" ADD CONSTRAINT " + foreignKey.Name +
		" FOREIGN KEY (" + foreignKey.Property + ")" +
		" REFERENCES " + foreignKey.ParentTable + " (" + foreignKey.References + ")" +
		" ON DELETE " + string(foreignKey.OnDelete) + ";"
}

// Drop generate SQL for dropping foreign key constraint
func (foreignKey *ForeignKey) Drop() string {
	return "ALTER TABLE " + foreignKey.TableName + " DROP CONSTRAINT " + foreignKey.Name + ";"
}
