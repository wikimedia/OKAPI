package upload

import (
	"okapi/helpers/bundle"
	"okapi/lib/task"
)

// Task for uploading the bundle
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	finish := func() error {
		return bundle.Generate(ctx.Project)
	}

	return nil, nil, finish
}
