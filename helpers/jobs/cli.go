package jobs

import (
	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/task"
)

// FromCLI get task instance and context
func FromCLI(params *cmd.Params) (job task.Task, ctx task.Context, err error) {
	job, err = getTask(*params.Task)

	if err != nil {
		return
	}

	project, err := getProject(*params.DBName)

	if err != nil {
		return
	}

	ctx = task.Context{
		State:   getState(*params.Task, *params.DBName),
		Project: &project,
		Log:     logger.Job,
		Params: task.Params{
			DBName:  *params.DBName,
			Restart: *params.Restart,
			Workers: *params.Workers,
			Limit:   *params.Limit,
			Offset:  *params.Offset,
			Pointer: *params.Pointer,
		},
	}

	return
}
