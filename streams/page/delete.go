package pageDelete

import (
	"encoding/json"
	"github.com/gookit/event"
	"github.com/r3labs/sse"
	"okapi/events/page"
)

// Payload event payload
type Payload struct {
	Database  string `json:"database"`
	PageTitle string `json:"page_title"`
}

// Handler page delete event handler for sse
func Handler(streamEvent *sse.Event) {
	payload := Payload{}

	json.Unmarshal(streamEvent.Data, &payload)

	event.Fire(pageDelete.Name, map[string]interface{}{
		"payload": pageDelete.Payload{
			Title:  payload.PageTitle,
			DBName: payload.Database,
		},
	})
}
