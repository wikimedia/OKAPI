package boot

import (
	"fmt"
	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/jobs/bundle"
	"okapi/jobs/scan"
	"okapi/jobs/sync"
	"okapi/lib/cache"
	"okapi/lib/cmd"
	"okapi/lib/runner"
	"okapi/lib/task"
	"okapi/models"
	"os"
	"os/signal"
	"time"
)

// Runner function to run tasks through redis
func Runner() {
	client := cache.NewClient()
	defer client.Close()
	cache.Client().Set("runner/status", "online", 0)
	defer cleanup()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cleanup()
			os.Exit(0)
		}
	}()

	jobs := map[task.Name]task.Task{
		"bundle": bundle.Task,
		"scan":   scan.Task,
		"sync":   sync.Task,
	}

	channel := client.Subscribe("runner").Channel()
	for message := range channel {
		command, err := runner.ParseCommand([]byte(message.Payload))

		if err != nil {
			continue
		}

		job, exists := jobs[task.Name(command.Task)]

		if !exists {
			logger.RUNNER.Error(logger.Message{
				ShortMessage: "Task " + command.Task + " not found!",
			})
			continue
		}

		project := models.Project{}
		err = models.DB().Model(&project).Where("db_name = ?", command.DBName).Select()

		if err != nil {
			logger.RUNNER.Error(logger.Message{
				ShortMessage: "Project " + project.DBName + " not found!",
				FullMessage:  err.Error(),
			})
			continue
		}

		executor := runner.Executor{
			Namespace: "runner/" + command.Task + "/" + command.DBName,
			Handler: func() error {
				restart := true
				workers := project.Schedule[command.Task].Workers
				err := task.Execute(job, &task.Context{
					State:   state.New(command.Task+"_"+command.DBName, 24*time.Hour),
					Project: &project,
					Cmd: &cmd.Params{
						Task:     &command.Task,
						Project:  &command.DBName,
						Position: cmd.Context.Position,
						Restart:  &restart,
						Workers:  &workers,
					},
				})

				if err != nil {
					logger.RUNNER.Error(logger.Message{
						ShortMessage: fmt.Sprintf("Task: task '%s' exec error", command.Task),
						FullMessage:  err.Error(),
					})
				}

				return err
			},
		}

		go executor.Run()
	}
}

func cleanup() {
	cache.Client().Del("runner/status")
}
