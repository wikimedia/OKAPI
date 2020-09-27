package revision

import (
	"github.com/gookit/event"
)

// Instance initialized instance of event
var Instance = &Event{}

// Name event name
const Name = "page_revision"

// Payload revision event payload
type Payload struct {
	Title    string
	Revision int
	DBName   string
	Redirect bool
	NsID     int
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
