package delete

import (
	"fmt"
	page_delete "okapi/events/page/delete"
	"okapi/lib/queue"
	"okapi/models"

	"github.com/gookit/event"
)

// Init add all event listeners
func Init() {
	event.On(page_delete.Name, event.ListenerFunc(Listener))
}

// Listener handler for event
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(page_delete.Payload)
	project := models.Project{}
	err := models.DB().Model(&project).Column("id").Where("db_name = ?", payload.DBName).Select()

	if err != nil {
		return err
	}

	page := models.Page{}

	models.DB().Model(&page).Where("title = ? and project_id = ?", payload.Title, project.ID).Select()

	if page.ID <= 0 {
		return fmt.Errorf("page does not exist: title: %q, project_id: %d", payload.Title, project.ID)
	}

	queue.PageDelete.Add(page)

	return nil
}
