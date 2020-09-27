package revision

import (
	page_revision "okapi/events/page/revision"
	"okapi/lib/queue"
	"okapi/models"
	page_revision_processor "okapi/processors/page/revision"

	"github.com/gookit/event"
)

// Init add all event listeners
func Init() {
	event.On(page_revision.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(page_revision.Payload)

	if payload.Redirect {
		return nil
	}

	project := models.Project{}
	page := models.Page{}

	err := models.DB().Model(&project).Where("db_name = ?", payload.DBName).Select()

	if err != nil {
		return err
	}

	err = models.DB().Model(&page).Where("title = ? and project_id = ?", payload.Title, project.ID).Select()

	if err != nil {
		page.Title = payload.Title
		page.Revision = 0
		page.SiteURL = project.SiteURL
		page.ProjectID = project.ID
		page.NsID = payload.NsID
		err = models.Save(&page)
	}

	if page.Revision != payload.Revision {
		project.Updates++
		models.Save(&project)
	}

	if err != nil {
		return err
	}

	queue.PageRevision.Add(page_revision_processor.Payload{
		Page:     page,
		Project:  project,
		Revision: payload.Revision,
	})

	return nil
}
