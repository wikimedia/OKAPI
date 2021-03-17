package models

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

// Project database table representation
type Project struct {
	ID           int       `pg:",pk" json:"id"`
	DbName       string    `pg:"type:varchar(255),unique,notnull" json:"db_name"`
	SiteName     string    `pg:"type:varchar(255),notnull" json:"site_name"`
	SiteCode     string    `pg:"type:varchar(255),notnull" json:"site_code"`
	SiteURL      string    `pg:"type:varchar(255),notnull" json:"site_url"`
	Lang         string    `pg:"type:varchar(25),notnull" json:"lang"`
	Active       bool      `pg:",use_zero,notnull" json:"active"`
	Threshold    float64   `pg:",notnull" json:"threshold"`
	Delay        int       `pg:",notnull" json:"delay"`
	HTMLSize     float64   `pg:"type:double precision" json:"html_size"`
	WikitextSize float64   `pg:"type:double precision" json:"wikitext_size"`
	JSONSize     float64   `pg:"type:double precision" json:"json_size"`
	HTMLPath     string    `pg:"type:varchar(255)" json:"html_path"`
	WikitextPath string    `pg:"type:varchar(255)" json:"wikitext_path"`
	JSONPath     string    `pg:"type:varchar(255)" json:"json_path"`
	HTMLAt       time.Time `json:"html_at"`
	WikitextAt   time.Time `json:"wikitext_at"`
	JSONAt       time.Time `json:"json_at"`
	Language     *Language `pg:"rel:has-one" json:"language,omitempty"`
	timestamp
}

var _ pg.BeforeUpdateHook = (*Project)(nil)

// BeforeUpdate model hook
func (proj *Project) BeforeUpdate(ctx context.Context) (context.Context, error) {
	proj.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Project)(nil)

// BeforeInsert model hook
func (proj *Project) BeforeInsert(ctx context.Context) (context.Context, error) {
	proj.OnInsert()
	return ctx, nil
}
