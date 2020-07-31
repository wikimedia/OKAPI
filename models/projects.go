package models

import (
	"context"
	"okapi/lib/schedule"
	"time"

	"github.com/go-pg/pg/v9"
)

// Project struct for "projects" table representation
type Project struct {
	baseModel
	Name      string                    `pg:"type:varchar(255),notnull" json:"name"`
	LocalName string                    `pg:"type:varchar(255),notnull" json:"local_name"`
	Code      string                    `pg:"type:varchar(25),notnull" json:"code"`
	SiteName  string                    `pg:"type:varchar(255),notnull" json:"site_name"`
	SiteCode  string                    `pg:"type:varchar(255),notnull" json:"site_code"`
	SiteURL   string                    `pg:"type:varchar(255),notnull" json:"site_url"`
	DBName    string                    `pg:"type:varchar(255),unique,notnull" json:"db_name"`
	Dir       string                    `pg:"type:varchar(3),notnull" json:"dir"`
	Size      float64                   `pg:"type:double precision" json:"size"`
	Path      string                    `pg:"type:varchar(255)" json:"path"`
	Active    bool                      `pg:",use_zero,notnull" json:"active"`
	Schedule  map[string]*schedule.Info `pg:"type:jsonb,notnull" json:"schedule"`
	Pages     []*Page                   `pg:"-" json:"pages"`
	DumpedAt  time.Time                 `json:"dumped_at"`
}

var _ pg.AfterSelectHook = (*Project)(nil)

// AfterSelect model hook
func (project *Project) AfterSelect(ctx context.Context) error {
	if project.Schedule == nil {
		project.Schedule = make(map[string]*schedule.Info)
	}

	for _, task := range []string{"scan", "sync", "bundle", "general"} {
		if _, ok := project.Schedule[task]; !ok {
			project.Schedule[task] = &schedule.Info{
				Workers:   10,
				Frequency: schedule.Daily,
			}
		}
	}

	return nil
}

var _ pg.BeforeUpdateHook = (*Project)(nil)

// BeforeUpdate model hook
func (project *Project) BeforeUpdate(ctx context.Context) (context.Context, error) {
	project.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Project)(nil)

// BeforeInsert model hook
func (project *Project) BeforeInsert(ctx context.Context) (context.Context, error) {
	project.OnInsert()
	return ctx, nil
}
