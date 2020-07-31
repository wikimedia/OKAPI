package revision

import (
	"github.com/gookit/event"
	"okapi/events/revision"
)

// Init add all event listeners
func Init() {
	name := revision.Name
	event.On(name, event.ListenerFunc(Page))
}
