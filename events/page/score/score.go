package score

import (
	"okapi/lib/ores"

	"github.com/gookit/event"
)

// Instance initialized instance of event
var Instance = &Event{}

// Name event name
const Name = "page_score"

// Payload revision event payload
type Payload struct {
	Title    string
	Revision int
	DBName   string
	Redirect bool
	NsID     int
	Scores   map[string]ores.Stream
}

// Event revision happened event
type Event struct {
	event.BasicEvent
}

// Init function to initialize the event
func Init() {
	Instance.SetName(Name)
	event.AddEvent(Instance)
}
