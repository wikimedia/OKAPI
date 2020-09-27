package boot

import (
	"fmt"

	"okapi/helpers/jobs"
	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/task"
)

// Task function to run task starter
func Task() {
	message := fmt.Sprintf("task '%s' exec error", *cmd.Context.Task)
	job, context, err := jobs.FromCLI(cmd.Context)

	if err != nil {
		logger.System.Panic(message, err.Error())
	}

	err = task.Exec(job, &context)

	if err != nil {
		logger.System.Panic(message, err.Error())
	}
}
