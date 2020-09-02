package bundle

import (
	"fmt"
	"okapi/helpers/logger"
	"okapi/lib/task"

	"gopkg.in/gookit/color.v1"
)

// Name task name for trigger
var Name task.Name = "bundle"

// Task bunnding files to s3 bucket
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	options := Options{}
	err := options.Init(ctx)

	if err == nil {
		go writeWorker(&options)
	}

	return Pool(ctx, &options), Worker(ctx, &options), Finish(ctx, &options), err
}

func writeWorker(options *Options) {
	for payload := range options.Writer.Queue {
		if err := options.Writer.Worker(payload); err != nil {
			logger.Job.Error(fmt.Sprintf("Worker #write - '%s'", payload.Name), err.Error(), map[string]interface{}{
				"_title": payload.Name,
			})
		} else {
			color.Yellow.Println("Worker #write - '" + payload.Name + "'")
		}
	}
}
