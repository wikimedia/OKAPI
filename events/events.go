package events

import (
	"okapi/events/page"
	"okapi/events/revision"
)

// Init initialize events
func Init() {
	events := []func(){
		revision.Init,
		pageDelete.Init,
	}

	for _, initializer := range events {
		initializer()
	}
}
