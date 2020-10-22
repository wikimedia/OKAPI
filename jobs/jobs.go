package jobs

import (
	"okapi/jobs/export"
	"okapi/jobs/index"
	"okapi/jobs/namespaces"
	"okapi/jobs/projects"
	"okapi/jobs/pull"
	"okapi/jobs/scan"
	"okapi/jobs/search"
	"okapi/jobs/upload"
	"okapi/lib/task"
)

// Tasks list of all available tasks
var Tasks = map[task.Name]task.Task{
	namespaces.Name: namespaces.Task,
	projects.Name:   projects.Task,
	search.Name:     search.Task,
	upload.Name:     upload.Task,
	export.Name:     export.Task,
	scan.Name:       scan.Task,
	pull.Name:       pull.Task,
	index.Name:      index.Task,
}
