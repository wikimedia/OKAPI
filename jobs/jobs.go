package jobs

import (
	"okapi/jobs/bundle"
	"okapi/jobs/export"
	"okapi/jobs/projects"
	"okapi/jobs/pull"
	"okapi/jobs/scan"
	"okapi/lib/task"
)

// Tasks list of all available tasks
var Tasks = map[task.Name]task.Task{
	projects.Name: projects.Task,
	export.Name:   export.Task,
	bundle.Name:   bundle.Task,
	scan.Name:     scan.Task,
	pull.Name:     pull.Task,
}
