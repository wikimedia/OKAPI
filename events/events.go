package events

import (
	page_delete "okapi/events/page/delete"
	page_revision "okapi/events/page/revision"
	page_score "okapi/events/page/score"
	page_unbundle "okapi/events/page/unbundle"
	project_bundle "okapi/events/project/bundle"
)

// Init initialize events
func Init() error {
	events := []func(){
		page_revision.Init,
		page_score.Init,
		page_delete.Init,
		page_unbundle.Init,
		project_bundle.Init,
	}

	for _, initializer := range events {
		initializer()
	}

	return nil
}
