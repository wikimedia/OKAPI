package jobs

import (
	"fmt"
	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/jobs"
	"okapi/lib/cmd"
	"okapi/lib/run"
	"okapi/lib/task"
	"okapi/models"
	"time"
)

func getState(name string, dbName string) *state.Client {
	return state.New(name+"_"+dbName, 24*time.Hour)
}

func getTask(name string) (job task.Task, err error) {
	job, exists := jobs.Tasks[task.Name(name)]

	if !exists {
		err = fmt.Errorf("task '%s' not found", string(name))
		return
	}

	return
}

func getProject(dbName string) (project models.Project, err error) {
	if dbName != "*" && len(dbName) > 0 {
		err = models.DB().Model(&project).Where("db_name = ?", dbName).Select()
	}

	return
}

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
		},
	}

	return
}

// FromCMD get job info from cmd
func FromCMD(cmd *run.Cmd, params *cmd.Params) (job task.Task, ctx task.Context, err error) {
	job, err = getTask(cmd.Task)

	if err != nil {
		return
	}

	project, err := getProject(cmd.DBName)

	if err != nil {
		return
	}

	workers := *params.Workers

	if schedule, ok := project.Schedule[cmd.Task]; ok && project.ID > 0 {
		workers = schedule.Workers
	}

	ctx = task.Context{
		State:   getState(cmd.Task, cmd.DBName),
		Project: &project,
		Log:     logger.Job,
		Params: task.Params{
			Restart: true,
			DBName:  cmd.DBName,
			Workers: workers,
			Limit:   *params.Limit,
			Offset:  *params.Offset,
		},
	}

	return
}
