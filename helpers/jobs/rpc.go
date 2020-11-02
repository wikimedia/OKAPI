package jobs

import (
	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/task"
	"okapi/protos/runner"
)

// FromRPC get job info from cmd
func FromRPC(req *runner.Request, params *cmd.Params) (job task.Task, ctx task.Context, err error) {
	job, err = getTask(req.Task)

	if err != nil {
		return
	}

	project, err := getProject(req.Database)

	if err != nil {
		return
	}

	workers := *params.Workers

	if schedule, ok := project.Schedule[req.Task]; ok && project.ID > 0 {
		workers = schedule.Workers
	}

	ctx = task.Context{
		State:   getState(req.Task, req.Database),
		Project: &project,
		Log:     logger.Job,
		Params: task.Params{
			Restart: true,
			DBName:  req.Database,
			Workers: workers,
			Limit:   *params.Limit,
			Offset:  *params.Offset,
		},
	}

	return
}
