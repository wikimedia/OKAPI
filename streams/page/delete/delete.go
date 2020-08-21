package delete

import (
	"encoding/json"
	"okapi/events/page/delete"
	page_delete "okapi/events/page/delete"

	"github.com/gookit/event"
	"github.com/r3labs/sse"
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

	event.Fire(page_delete.Name, map[string]interface{}{
		"payload": delete.Payload{
			Title:  payload.PageTitle,
			DBName: payload.Database,
		},
	})
}
