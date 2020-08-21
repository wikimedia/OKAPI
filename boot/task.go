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
		logger.SYSTEM.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Task: task '%s' not found", string(name)),
		})
	}

	if *cmd.Context.Project != "*" {
		models.DB().Model(&project).Where("db_name = ?", *cmd.Context.Project).Select()
	}

	if *cmd.Context.Project != "*" && project.ID <= 0 {
		logger.SYSTEM.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Task: project '%s' not found", *cmd.Context.Project),
		})
	}

	err := task.Execute(job, &task.Context{
		Cmd:     cmd.Context,
		State:   state.New(*cmd.Context.Task+"_"+*cmd.Context.Project, 24*time.Hour),
		Project: &project,
	})

	if err != nil {
		logger.SYSTEM.Panic(logger.Message{
			ShortMessage: fmt.Sprintf("Task: task '%s' exec error", string(name)),
			FullMessage:  err.Error(),
		})
	}
}
