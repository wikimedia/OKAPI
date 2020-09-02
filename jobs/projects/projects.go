package projects

import (
	"okapi/lib/task"
	"okapi/lib/wiki"
)

// Name task name for trigger
var Name task.Name = "projects"

// Task getting all projects from sitematrix
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	projects, _, err := wiki.Client("https://en.wikipedia.org/").GetSitematrix()
	return Pool(projects), Worker, nil, err
}
