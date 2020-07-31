package scan

import (
	"fmt"

	"okapi/helpers/logger"
	"okapi/lib/task"
)

// Task to get titles from dump and store them in database
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	options := Options{}
	worker := func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message, info, err := Worker(id, payload)
		options.Position++
		ctx.State.Set("position", options.Position)
		return message, info, err
	}

	err := options.Init(ctx)
	if err != nil {
		logger.JOB.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Job: 'bundle' for the project '%s' exec failed", *ctx.Cmd.Project),
			FullMessage:  err.Error(),
		})
	}

	return Pool(ctx, &options), worker, nil
}
