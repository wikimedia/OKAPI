package revision

import (
	"encoding/json"

	page_revision "okapi/events/page/revision"

	"github.com/gookit/event"
	"github.com/r3labs/sse"
)

// Payload event payload
type Payload struct {
	PageTitle      string `json:"page_title"`
	Database       string `json:"database"`
	RevID          int    `json:"rev_id"`
	PageIsRedirect bool   `json:"page_is_redirect"`
	PageNamespace  int    `json:"page_namespace"`
}

// Handler revision event handler for sse
func Handler(streamEvent *sse.Event) {
	payload := Payload{}
	json.Unmarshal(streamEvent.Data, &payload)
	event.Fire(page_revision.Name, map[string]interface{}{
		"payload": page_revision.Payload{
			Title:    payload.PageTitle,
			Revision: payload.RevID,
			DBName:   payload.Database,
			Redirect: payload.PageIsRedirect,
			NsID:     payload.PageNamespace,
		},
	})
}
