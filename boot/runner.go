package boot

import (
	"fmt"
	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/jobs"
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

	channel := client.Subscribe(runner.Channel).Channel()
	for message := range channel {
		command, err := runner.ParseCommand([]byte(message.Payload))
		namespace := "runner/" + command.Task + "/" + command.DBName

		if err != nil {
			continue
		}

		job, exists := jobs.Tasks[task.Name(command.Task)]

		if !exists {
			message := "Task " + command.Task + " not found!"

			logger.RUNNER.Error(logger.Message{
				ShortMessage: message,
			})

			runner.Error.Send(namespace, &runner.Message{
				Info: message,
			})
			continue
		}

		project := models.Project{}
		err = models.DB().Model(&project).Where("db_name = ?", command.DBName).Select()

		if err != nil {
			message := "Project " + project.DBName + " not found!"

			logger.RUNNER.Error(logger.Message{
				ShortMessage: message,
				FullMessage:  err.Error(),
			})

			runner.Error.Send(namespace, &runner.Message{
				Info: message,
			})
			continue
		}

		executor := runner.Executor{
			Namespace: namespace,
			Handler: func() error {
				restart := true
				workers := *cmd.Context.Workers

				if schedule, ok := project.Schedule[command.Task]; ok {
					workers = schedule.Workers
				}

				err := task.Execute(job, &task.Context{
					State:   state.New(command.Task+"_"+command.DBName, 24*time.Hour),
					Project: &project,
					Cmd: &cmd.Params{
						Task:    &command.Task,
						Project: &command.DBName,
						Offset:  cmd.Context.Offset,
						Limit:   cmd.Context.Limit,
						Restart: &restart,
						Workers: &workers,
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
	cache.Client().Del(runner.Channel + "/status")
}
