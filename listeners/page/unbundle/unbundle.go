package unbundle

import (
	page_unbundle "okapi/events/page/unbundle"
	"okapi/helpers/damaging"
	"okapi/jobs/export"
	"okapi/lib/run"

	"github.com/gookit/event"
)

// Init add all event listeners
func Init() {
	event.On(page_unbundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(page_unbundle.Payload)
	err := damaging.Add(payload.Title, payload.DBName)

	if err != nil {
		return err
	}

	cmd := run.Cmd{
		Task:   string(export.Name),
		DBName: payload.DBName,
	}

	return cmd.Enqueue()
}
