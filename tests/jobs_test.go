package main

import (
	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/jobs/export"
	"okapi/jobs/pull"
	"okapi/jobs/scan"
	"okapi/lib/cache"
	"okapi/lib/cmd"
	"okapi/lib/env"
	"okapi/lib/storage"
	"okapi/lib/task"
	"okapi/models"
	"testing"
	"time"
)

func TestJobs(t *testing.T) {
	cmd.Context.Parse()
	env.Context.Parse(".env")
	storage.Init()
	models.Init()
	cache.Init()
	defer models.Close()
	defer cache.Close()

	project := models.Project{
		LangName:      "English",
		LangLocalName: "English",
		Lang:          "en",
		SiteName:      "Wikipedia",
		SiteCode:      "wiki",
		SiteURL:       "https://en.wikipedia.org",
		DBName:        "enwiki_test",
		Dir:           "ltr",
		Active:        true,
	}

	models.DB().Model(&project).Where("db_name = ?", project.DBName).SelectOrInsert()
	err := models.DB().Model(&project).Where("db_name = ?", project.DBName).Select()

	if err != nil {
		t.Error("Wiki `" + project.DBName + "` not found -> " + err.Error())
		return
	}

	params := task.Params{
		DBName:  project.DBName,
		Restart: true,
		Workers: 5,
		Limit:   100,
		Offset:  0,
	}

	tasks := []func() error{
		func() error {
			return task.Exec(scan.Task, &task.Context{
				Params:  params,
				State:   state.New("scan_"+project.DBName, 24*time.Hour),
				Project: &project,
				Log:     logger.Job,
			})
		},
		func() error {
			return task.Exec(pull.Task, &task.Context{
				Params:  params,
				State:   state.New("pull_"+project.DBName, 24*time.Hour),
				Project: &project,
				Log:     logger.Job,
			})
		},
		func() error {
			return task.Exec(export.Task, &task.Context{
				Params:  params,
				State:   state.New("export_"+project.DBName, 24*time.Hour),
				Project: &project,
				Log:     logger.Job,
			})
		},
	}

	for _, task := range tasks {
		err = task()
		if err != nil {
			t.Error(err)
			return
		}
	}

	_, err = storage.Local.Client().Get(project.Path)

	if err != nil {
		t.Error(err)
	}
}
