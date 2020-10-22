package models

import (
	"context"

	"github.com/go-pg/pg/v10"
)

// Export struct for "exports" table representation
type Export struct {
	timestampModel
	UserID       int    `pg:",notnull,pk" json:"user_id"`
	ResourceType string `pg:"type:varchar(25),notnull,pk" json:"resource_type" binding:"required,export_type"`
	ResourceID   int    `pg:",notnull,pk" json:"resource_id"`
}

// IsUpdate check if it's update or insert
func (export *Export) IsUpdate() bool {
	isUpdate, _ := db.Model(export).Where(
		"user_id = ? and resource_type = ? and resource_id = ?",
		export.ResourceType, export.ResourceType, export.ResourceID,
	).Exists()
	return isUpdate
}

var _ pg.BeforeUpdateHook = (*Export)(nil)

// BeforeUpdate model hook
func (export *Export) BeforeUpdate(ctx context.Context) (context.Context, error) {
	export.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Export)(nil)

// BeforeInsert model hook
func (export *Export) BeforeInsert(ctx context.Context) (context.Context, error) {
	export.OnInsert()
	return ctx, nil
}
