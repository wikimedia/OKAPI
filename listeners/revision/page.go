package revision

import (
	"github.com/gookit/event"
	"okapi/events/revision"
	"okapi/lib/queue"
	"okapi/models"
)

// Page event to update page on revision
func Page(e event.Event) error {
	payload := e.Data()["payload"].(revision.Payload)
	project := models.Project{}
	err := models.DB().Model(&project).Column("id", "site_url", "code").Where("db_name = ?", payload.DBName).Select()

	if err != nil {
		return err
	}

	page := models.Page{}
	models.DB().Model(&page).Where("title = ? and project_id = ?", payload.Title, project.ID).Select()

	if page.ID <= 0 {
		page.Revision = payload.Revision
		page.Title = payload.Title
		page.SiteURL = project.SiteURL
		page.ProjectID = project.ID
		page.Lang = project.Code
	}

	queue.Scan.Add(page)

	return nil
}
