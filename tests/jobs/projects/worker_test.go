package projects

import (
	"okapi/jobs/projects"
	"okapi/lib/db"
	"okapi/lib/env"
	"okapi/models"
	"testing"
)

func TestWorker(t *testing.T) {
	env.Context.Fill()
	db := db.Client()
	defer db.Close()

	project := models.Project{
		Name:      "English",
		Code:      "en",
		SiteName:  "Wikipedia",
		SiteURL:   "https://en.wikipedia.org",
		SiteCode:  "wiki",
		DBName:    "enwiki_new",
		Dir:       "ltr",
		LocalName: "English",
		Active:    true,
	}

	_, _, err := projects.Worker(1, &project)

	if err != nil {
		t.Error(err)
	}

	if project.ID <= 0 {
		t.Error("Project `" + project.Name + "` wasn't created")
	}
}
