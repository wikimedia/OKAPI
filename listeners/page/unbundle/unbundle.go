package unbundle

import (
	page_unbundle "okapi/events/page/unbundle"
	"okapi/helpers/damaging"
	"okapi/lib/runner"

	"github.com/gookit/event"
)

// Init add all event listeners
func Init() {
	event.On(page_unbundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(page_unbundle.Payload)
	err := damaging.Add(payload.Revision, payload.DBName)

	if err != nil {
		return err
	}

	command := runner.Command{
		Task:   "bundle",
		DBName: payload.DBName,
	}

	return command.Exec()
}
