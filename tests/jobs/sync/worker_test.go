package sync

import (
	"okapi/jobs/sync"
	"okapi/lib/db"
	"okapi/lib/env"
	"okapi/lib/storage"
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
		DBName:    "enwiki_sync",
		Dir:       "ltr",
		LocalName: "English",
		Active:    true,
	}

	models.DB().Model(&project).Where("db_name = ?", project.DBName).SelectOrInsert()

	page := models.Page{
		Title:     "Okapi",
		ProjectID: project.ID,
		Revision:  1234,
		TID:       "tid",
		Lang:      project.Code,
		SiteURL:   project.SiteURL,
	}

	models.DB().Model(&page).
		Where("title = ? and project_id = ?", page.Title, page.ProjectID).
		SelectOrInsert()

	_, _, err := sync.Worker(1, &page)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = storage.Local.Client().Get(page.Path)

	if err != nil {
		t.Error(err)
		return
	}
}
