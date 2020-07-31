package pageDelete

import (
	"fmt"
	"github.com/gookit/event"
	"okapi/events/page"
	"okapi/lib/queue"
	"okapi/models"
)

func Page(e event.Event) error {
	payload := e.Data()["payload"].(pageDelete.Payload)
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

	queue.DeletePage.Add(page)

	return nil
}
