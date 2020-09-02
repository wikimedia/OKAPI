package unbundle

import (
	"github.com/gookit/event"
)

// Instance initialized instance of event
var Instance = &Event{}

// Name event name
const Name = "page_unbundle"

// Payload revision event payload
type Payload struct {
	Revision int
	Title    string
	DBName   string
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
