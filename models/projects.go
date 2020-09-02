package models

import (
	"context"
	"okapi/lib/env"
	"okapi/lib/ores"
	"okapi/lib/schedule"
	"time"

	"github.com/go-pg/pg/v9"
)

// Project struct for "projects" table representation
type Project struct {
	baseModel
	SiteName      string                    `pg:"type:varchar(255),notnull" json:"site_name"`
	SiteCode      string                    `pg:"type:varchar(255),notnull" json:"site_code"`
	SiteURL       string                    `pg:"type:varchar(255),notnull" json:"site_url"`
	DBName        string                    `pg:"type:varchar(255),unique,notnull" json:"db_name"`
	Lang          string                    `pg:"type:varchar(25),notnull" json:"lang"`
	LangName      string                    `pg:"type:varchar(255),notnull" json:"lang_name"`
	LangLocalName string                    `pg:"type:varchar(255),notnull" json:"lang_local_name"`
	Dir           string                    `pg:"type:varchar(3),notnull" json:"dir"`
	Size          float64                   `pg:"type:double precision" json:"size"`
	Path          string                    `pg:"type:varchar(255)" json:"path"`
	Active        bool                      `pg:",use_zero,notnull" json:"active"`
	Schedule      map[string]*schedule.Info `pg:"type:jsonb,notnull" json:"schedule"`
	Threshold     map[ores.Model]float64    `pg:",notnull" json:"threshold" binding:"required,threshold"`
	TimeDelay     int                       `pg:",use_zero" json:"time_delay" binding:"required,number"` // in hours
	Updates       int                       `pg:",use_zero" json:"updates"`
	Pages         []*Page                   `pg:"-" json:"pages"`
	DumpedAt      time.Time                 `json:"dumped_at"`
}

// GetExportName get name of the export
func (project *Project) GetExportName() string {
	return "export_" + project.DBName + ".tar.gz"
}

// GetExportPath get path of the export
func (project *Project) GetExportPath() string {
	return env.Context.VolumeMountPath + "/exports/" + project.DBName + "/" + project.GetExportName()
}

// GetRemoteExportPath get remote path of the export
func (project *Project) GetRemoteExportPath() string {
	return "/exports/" + project.DBName + "/" + project.GetExportName()
}

// GetThreshold getting the value by name
func (project *Project) GetThreshold(oresModel ores.Model) *float64 {
	if project.Threshold != nil {
		if val, exists := project.Threshold[oresModel]; exists {
			return &val
		}
	}

	return nil
}

var _ pg.AfterSelectHook = (*Project)(nil)

// AfterSelect model hook
func (project *Project) AfterSelect(ctx context.Context) error {
	if project.Schedule == nil {
		project.Schedule = make(map[string]*schedule.Info)
	}

	for _, task := range schedule.Jobs {
		if _, ok := project.Schedule[task]; !ok {
			project.Schedule[task] = &schedule.Info{
				Workers:   250,
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
