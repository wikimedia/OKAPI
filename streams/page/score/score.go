package score

import (
	"encoding/json"
	"okapi/lib/ores"

	page_score "okapi/events/page/score"

	"github.com/gookit/event"
	"github.com/r3labs/sse"
)

// Payload event payload
type Payload struct {
	PageTitle      string                 `json:"page_title"`
	Database       string                 `json:"database"`
	RevID          int                    `json:"rev_id"`
	PageIsRedirect bool                   `json:"page_is_redirect"`
	PageNamespace  int                    `json:"page_namespace"`
	Scores         map[string]ores.Stream `json:"scores"`
}

// Handler revision event handler for sse
func Handler(streamEvent *sse.Event) {
	payload := Payload{}
	json.Unmarshal(streamEvent.Data, &payload)

	event.Fire(page_score.Name, map[string]interface{}{
		"payload": page_score.Payload{
			Title:    payload.PageTitle,
			Revision: payload.RevID,
			DBName:   payload.Database,
			Redirect: payload.PageIsRedirect,
			NsID:     payload.PageNamespace,
			Scores:   payload.Scores,
		},
	})
}
