package main

import (
	"okapi/helpers/state"
	"okapi/jobs/bundle"
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
	storage.Local.Client()
	storage.Remote.Client()
	db := models.DB()
	ch := cache.Client()
	defer db.Close()
	defer ch.Close()

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
	cmd.Context.Project = &project.DBName
	err := models.DB().Model(&project).Where("db_name = ?", project.DBName).Select()

	if err != nil {
		t.Error("Wiki `" + project.DBName + "` not found -> " + err.Error())
		return
	}

	tasks := []func() error{
		func() error {
			return task.Execute(scan.Task, &task.Context{
				Cmd:     cmd.Context,
				State:   state.New("scan_"+project.DBName, 24*time.Hour),
				Project: &project,
			})
		},
		func() error {
			return task.Execute(pull.Task, &task.Context{
				Cmd:     cmd.Context,
				State:   state.New("pull_"+project.DBName, 24*time.Hour),
				Project: &project,
			})
		},
		func() error {
			return task.Execute(bundle.Task, &task.Context{
				Cmd:     cmd.Context,
				State:   state.New("bundle_"+project.DBName, 24*time.Hour),
				Project: &project,
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

	_, err = storage.Local.Client().Get(project.RelativeBundlePath())

	if err != nil {
		t.Error(err)
	}
}
