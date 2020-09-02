package boot

import (
	"fmt"
	"time"

	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/jobs"
	"okapi/lib/cmd"
	"okapi/lib/task"
	"okapi/models"
)

// Task function to run task starter
func Task() {
	name := task.Name(*cmd.Context.Task)
	project := models.Project{}
	job, exists := jobs.Tasks[name]

	if !exists {
		logger.System.Panic("Task startup failed", fmt.Sprintf("Task '%s' not found", string(name)))
	}

	if *cmd.Context.Project != "*" {
		models.DB().Model(&project).Where("db_name = ?", *cmd.Context.Project).Select()
	}

	if *cmd.Context.Project != "*" && project.ID <= 0 {
		logger.System.Panic("Task startup failed", fmt.Sprintf("Project '%s' not found", *cmd.Context.Project))
	}

	err := task.Execute(job, &task.Context{
		Cmd:     cmd.Context,
		State:   state.New(*cmd.Context.Task+"_"+*cmd.Context.Project, 24*time.Hour),
		Project: &project,
	})

	if err != nil {
		logger.System.Panic(fmt.Sprintf("Task '%s' exec error", string(name)), err.Error())
	}
}
