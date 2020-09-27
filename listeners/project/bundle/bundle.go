package bundle

import (
	project_bundle "okapi/events/project/bundle"
	"okapi/jobs/export"
	"okapi/lib/run"

	"github.com/gookit/event"
)

// Init setup a listener
func Init() {
	event.On(project_bundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(project_bundle.Payload)

	cmd := run.Cmd{
		Task:   string(export.Name),
		DBName: payload.DBName,
	}

	return cmd.Enqueue()
}
