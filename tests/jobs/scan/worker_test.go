package scan

import (
	"okapi/jobs/scan"
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
		DBName:    "enwiki_scan",
		Dir:       "ltr",
		LocalName: "English",
		Active:    true,
	}

	models.DB().Model(&project).SelectOrInsert()

	page := models.Page{
		Title:     "United_States",
		SiteURL:   "https://en.wikipedia.org",
		ProjectID: project.ID,
	}
	_, _, err := scan.Worker(1, &page)

	if err != nil {
		t.Error(err)
		return
	}

	if page.ID <= 0 || page.Revision <= 0 {
		t.Error("Page `" + page.Title + "` wasn't updated in the database.")
	}
}
