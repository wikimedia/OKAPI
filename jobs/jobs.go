package jobs

import (
	"okapi/jobs/bundle"
	"okapi/jobs/projects"
	"okapi/jobs/pull"
	"okapi/jobs/scan"
	"okapi/jobs/upload"
	"okapi/lib/task"
)

// Tasks list of all available tasks
var Tasks = map[task.Name]task.Task{
	"projects": projects.Task,
	"bundle":   bundle.Task,
	"scan":     scan.Task,
	"pull":     pull.Task,
	"upload":   upload.Task,
}
