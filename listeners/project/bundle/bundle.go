package bundle

import (
	project_bundle "okapi/events/project/bundle"
	"okapi/lib/runner"

	"github.com/gookit/event"
)

// Init setup a listener
func Init() {
	event.On(project_bundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(project_bundle.Payload)

	command := runner.Command{
		Task:   "bundle",
		DBName: payload.DBName,
	}

	return command.Exec()
}
