package pageDelete

import (
	"github.com/gookit/event"
)

// Instance initialized instance of event
var Instance = &Event{}

// Name event name
const Name = "page-delete"

// Payload page delete event payload
type Payload struct {
	Title  string
	DBName string
}

// Event page delete happened event
type Event struct {
	event.BasicEvent
}

// Init function to initialize the event
func Init() {
	Instance.SetName(Name)
	event.AddEvent(Instance)
}
