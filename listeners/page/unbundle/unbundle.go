package unbundle

import (
	page_unbundle "okapi/events/page/unbundle"
	"okapi/helpers/bundle"
	"okapi/lib/runner"
	"okapi/models"

	"github.com/gookit/event"
)

// Init add all event listeners
func Init() {
	event.On(page_unbundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(page_unbundle.Payload)
	page := models.Page{}

	err := models.DB().
		Model(&page).
		Relation("Project").
		Where("title = ? and db_name = ?", payload.Title, payload.DBName).
		Select()

	if err != nil {
		return err
	}

	if err = bundle.Delete(page.Project, &page); err != nil {
		return err
	}

	command := runner.Command{
		Task:   "upload",
		DBName: page.Project.DBName,
	}

	return command.Exec()
}
