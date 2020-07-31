package listeners

import "okapi/listeners/revision"
import "okapi/listeners/page/delete"

// Init initialize event listeners
func Init() {
	listeners := []func(){
		revision.Init,
		pageDelete.Init,
	}

	for _, listener := range listeners {
		listener()
	}
}
