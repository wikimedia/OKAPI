package listeners

import (
	page_delete "okapi/listeners/page/delete"
	page_revision "okapi/listeners/page/revision"
	page_unbundle "okapi/listeners/page/unbundle"
	project_bundle "okapi/listeners/project/bundle"
)

// Init initialize event listeners
func Init() {
	listeners := []func(){
		page_delete.Init,
		page_revision.Init,
		page_unbundle.Init,
		project_bundle.Init,
	}

	for _, listener := range listeners {
		listener()
	}
}
