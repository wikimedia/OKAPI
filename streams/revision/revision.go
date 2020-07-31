package revision

import (
	"encoding/json"

	"github.com/gookit/event"
	"github.com/r3labs/sse"
	"okapi/events/revision"
)

// Payload event payload
type Payload struct {
	PageTitle string `json:"page_title"`
	Database  string `json:"database"`
	RevID     int    `json:"rev_id"`
}

// Handler revision event handler for sse
func Handler(streamEvent *sse.Event) {
	payload := Payload{}
	json.Unmarshal(streamEvent.Data, &payload)
	event.Fire(revision.Name, map[string]interface{}{
		"payload": revision.Payload{
			Title:    payload.PageTitle,
			Revision: payload.RevID,
			DBName:   payload.Database,
		},
	})
}
