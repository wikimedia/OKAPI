package bundle

import (
	"fmt"

	"okapi/helpers/logger"
	"okapi/helpers/writer"
	"okapi/lib/task"
	"gopkg.in/gookit/color.v1"
)

// Task bunnding files to s3 bucket
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	options := Options{}
	err := options.Init(ctx)
	if err != nil {
		logger.JOB.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Job: 'bundle' for the project '%s' exec failed", *ctx.Cmd.Project),
			FullMessage:  err.Error(),
		})
	}

	go (func(queue <-chan writer.Payload) {
		message := "Worker #write, processed the unit! Message: page title: '%s'"
		for payload := range queue {
			if err := options.Writer.Worker(payload); err != nil {
				logger.JOB.Error(logger.Message{
					ShortMessage: fmt.Sprintf(message+", error: %s;", payload.Name, err),
					FullMessage:  err.Error(),
					Params: map[string]interface{}{
						"_title": payload.Name,
					},
				})
			} else {
				color.Yellow.Println(fmt.Sprintf(message+";", payload.Name))
			}
		}
	})(options.Writer.Queue)

	return Pool(ctx, &options), Worker(ctx, &options), Finish(ctx, &options)
}
