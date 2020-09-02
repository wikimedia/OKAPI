package runner

import (
	"fmt"
	"okapi/helpers/state"
	"okapi/lib/cache"
	lib_cmd "okapi/lib/cmd"
	"okapi/lib/task"
	"okapi/models"
	"time"
)

// Channel runner cache channel name
const Channel string = "runner"

// Online indicate that runner is online
func Online() error {
	return cache.Client().Set(Channel+"/status", "online", 0).Err()
}

// Offline indicate that runner is offline
func Offline() error {
	return cache.Client().Del(Channel + "/status").Err()
}

// IsOnline check if runner is online
func IsOnline() bool {
	return cache.Client().Exists(Channel+"/status").Val() == 1
}

// Execute the command for project
func Execute(cmd *Command) error {
	job, err := cmd.Job()

	if err != nil {
		return err
	}

	project, err := getProject(cmd)

	if err != nil {
		return err
	}

	if cmd.IsRunning() {
		return fmt.Errorf("Task '%s' for db_name '%s' already running", cmd.Task, cmd.DBName)
	}

	cmd.Running()
	defer cmd.Stop()

	err = task.Execute(job, getContext(cmd, project))
	if err != nil {
		return err
	}

	return nil
}

func getProject(cmd *Command) (*models.Project, error) {
	project := models.Project{}

	if len(cmd.DBName) <= 0 {
		return &project, nil
	}

	err := models.DB().Model(&project).Where("db_name = ?", cmd.DBName).Select()
	return &project, err
}

func getContext(cmd *Command, project *models.Project) *task.Context {
	workers := *lib_cmd.Context.Workers

	if schedule, ok := project.Schedule[cmd.Task]; ok {
		workers = schedule.Workers
	}

	return &task.Context{
		State:   state.New(cmd.Task+"_"+cmd.DBName, 24*time.Hour),
		Project: project,
		Cmd: &lib_cmd.Params{
			Task:    &cmd.Task,
			Project: &cmd.DBName,
			Offset:  lib_cmd.Context.Offset,
			Limit:   lib_cmd.Context.Limit,
			Restart: lib_cmd.Context.Restart,
			Workers: &workers,
		},
	}
}
