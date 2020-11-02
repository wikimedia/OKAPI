package bundle

import (
	"context"
	project_bundle "okapi/events/project/bundle"
	"okapi/jobs/export"
	"okapi/lib/run"
	"okapi/protos/runner"

	"github.com/gookit/event"
)

// Init setup a listener
func Init() {
	event.On(project_bundle.Name, event.ListenerFunc(Listener))
}

// Listener event handler
func Listener(e event.Event) error {
	payload := e.Data()["payload"].(project_bundle.Payload)
	client, err := run.Client()

	if err != nil {
		return err
	}

	_, err = client.Enqueue(context.Background(), &runner.Request{
		Task:     string(export.Name),
		Database: payload.DBName,
	})

	return err
}
