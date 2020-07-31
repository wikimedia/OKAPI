package pageDelete

import (
	"github.com/gookit/event"
	"okapi/events/page"
)

// Init add all event listeners
func Init() {
	name := pageDelete.Name
	event.On(name, event.ListenerFunc(Page))
}
