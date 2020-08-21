package pull

import (
	"okapi/jobs/pull"
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
		Lang:          "en",
		LangLocalName: "English",
		LangName:      "English",
		SiteName:      "Wikipedia",
		SiteURL:       "https://en.wikipedia.org",
		SiteCode:      "wiki",
		DBName:        "enwiki_pull",
		Dir:           "ltr",
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
		Title:     "Okapi",
		ProjectID: project.ID,
		SiteURL:   project.SiteURL,
	}

	page.SetRevision(973564695)

	_, err = models.DB().Model(&page).
		Where("title = ? and project_id = ?", page.Title, page.ProjectID).
		SelectOrInsert()

	if err != nil {
		t.Error(err)
		return
	}

	_, _, err = pull.Worker(1, &page)

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
