package bundle

import (
	"github.com/gookit/event"
)

// Instance initialized instance of event
var Instance = &Event{}

// Name event name
const Name = "project_bundle"

// Payload bundle event payload
type Payload struct {
	DBName string
}

// Event bundle struct
type Event struct {
	event.BasicEvent
}

// Init function to initialize the event
func Init() {
	Instance.SetName(Name)
	event.AddEvent(Instance)
}
