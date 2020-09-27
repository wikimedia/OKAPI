package scan

import (
	"okapi/jobs/scan"
	"okapi/lib/env"
	"okapi/models"
	"testing"
)

func TestWorker(t *testing.T) {
	env.Context.Fill()
	models.Init()
	defer models.Close()

	project := models.Project{
		Lang:          "en",
		LangName:      "English",
		LangLocalName: "English",
		SiteName:      "Wikipedia",
		SiteURL:       "https://en.wikipedia.org",
		SiteCode:      "wiki",
		DBName:        "enwiki_scan",
		Dir:           "ltr",
		TimeDelay:     0,
		Active:        true,
	}

	_, err := models.DB().
		Model(&project).
		Where("db_name = ?", project.DBName).
		SelectOrInsert()

	if err != nil {
		t.Error(err)
		return
	}

	page := models.Page{
		Title:     "United_States",
		SiteURL:   "https://en.wikipedia.org",
		Revision:  1,
		ProjectID: project.ID,
		Project:   &project,
	}
	_, _, err = scan.Worker(1, &page)

	if err != nil {
		t.Error(err)
		return
	}

	if page.ID <= 0 || page.Revisions[0] <= 0 {
		t.Error("Page `" + page.Title + "` wasn't updated in the database.")
	}
}
