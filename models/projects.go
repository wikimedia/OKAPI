package models

import (
	"context"
	"fmt"
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
	Threshold     map[ores.Model]float64    `pg:",notnull" json:"threshold"`
	TimeDelay     int                       `pg:",use_zero" json:"time_delay"` // in hours
	Updates       int                       `pg:",use_zero" json:"updates"`
	Pages         []*Page                   `pg:"-" json:"pages"`
	DumpedAt      time.Time                 `json:"dumped_at"`
}

// BundleName get file name for a  bundle
func (project *Project) BundleName() string {
	return "export_" + project.DBName + ".tar"
}

// RelativeBundlePath create bundle relative path
func (project *Project) RelativeBundlePath() string {
	return "/exports/" + project.DBName + "/" + project.BundleName()
}

// BundlePath get bundle path
func (project *Project) BundlePath() string {
	return env.Context.VolumeMountPath + project.RelativeBundlePath()
}

// RemoteBundlePath get bundle path for remote storage
func (project *Project) RemoteBundlePath() string {
	return project.DBName + "/" + project.CompressedBundleName()
}

// CompressedBundleName get compressed bundle name
func (project *Project) CompressedBundleName() string {
	return project.BundleName() + ".bz2"
}

// CompressedBundlePath get compressed bundle path
func (project *Project) CompressedBundlePath() string {
	return project.BundlePath() + ".bz2"
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

var thresholdModels = []ores.Model{ores.Damaging}

// thresholdValidation validates threshold model names
func (project *Project) thresholdValidation() error {
	if len(project.Threshold) == 0 {
		return nil
	}

	for modelName := range project.Threshold {
		isValidName := false

		for i := 0; i < len(thresholdModels); i++ {
			if modelName == thresholdModels[i] {
				isValidName = true

				break
			}
		}

		if !isValidName {
			return fmt.Errorf("\"%s\" threshold model name is not valid", modelName)
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

	for _, task := range []string{"scan", "pull", "bundle", "general"} {
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
	if err := project.thresholdValidation(); err != nil {
		return ctx, err
	}

	project.OnUpdate()
	return ctx, nil
}

var _ pg.BeforeInsertHook = (*Project)(nil)

// BeforeInsert model hook
func (project *Project) BeforeInsert(ctx context.Context) (context.Context, error) {
	if err := project.thresholdValidation(); err != nil {
		return ctx, err
	}

	project.OnInsert()
	return ctx, nil
}
