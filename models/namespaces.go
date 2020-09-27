package models

import (
	"context"

	"github.com/go-pg/pg/v9"
)

// Namespace wiki specific namespace
type Namespace struct {
	baseModel
	ID    int    `pg:",pk,use_zero" json:"id"`
	Title string `pg:"type:varchar(255),pk" json:"title"`
	Lang  string `pg:"type:varchar(25),notnull" json:"lang"`
}

// IsUpdate check if it's update or insert
func (ns *Namespace) IsUpdate() bool {
	isUpdate, _ := db.Model(ns).Where("id = ? and lang = ?", ns.ID, ns.Lang).Exists()
	return isUpdate
}

var _ pg.BeforeUpdateHook = (*Namespace)(nil)

// BeforeUpdate model hook
func (ns *Namespace) BeforeUpdate(ctx context.Context) (context.Context, error) {
	ns.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Namespace)(nil)

// BeforeInsert model hook
func (ns *Namespace) BeforeInsert(ctx context.Context) (context.Context, error) {
	ns.OnInsert()
	return ctx, nil
}
