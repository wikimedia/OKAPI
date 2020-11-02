package unbundle

import (
	"context"
	page_unbundle "okapi/events/page/unbundle"
	"okapi/helpers/damaging"
	"okapi/jobs/export"
	"okapi/lib/run"
	"okapi/protos/runner"

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
