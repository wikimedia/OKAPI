package pull

import (
	"fmt"
	"okapi/helpers/logger"
	"okapi/lib/task"
)

// Task for getting html for pages
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	options := Options{}

	err := options.Init(ctx)

	if err != nil {
		logger.JOB.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Job: 'pull' for the project '%s' exec failed", *ctx.Cmd.Project),
			FullMessage:  err.Error(),
		})
	}

	worker := func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message, info, err := Worker(id, payload)
		options.Position++
		ctx.State.Set("offset", options.Position)
		return message, info, err
	}

	return Pool(ctx, &options), worker, nil
}
