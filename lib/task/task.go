package task

import (
	"fmt"
	"sync"

	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/lib/cmd"
	"okapi/models"
)

// Context task execution context
type Context struct {
	Cmd     *cmd.Params
	State   *state.Client
	Project *models.Project
}

// Task struct
type Task func(ctx *Context) (Pool, Worker, Finish)

// Name task name
type Name string

// Payload for single worker
type Payload interface{}

// Worker single unit processor
type Worker func(id int, payload Payload) (string, map[string]interface{}, error)

// Pool function to get payload into the queue
type Pool func() ([]Payload, error)

// Finish function to run something after job finished
type Finish func() error

// Execute some task with context
func Execute(cmd Task, ctx *Context) error {
	if *ctx.Cmd.Restart {
		ctx.State.Clear()
	}

	wg := &sync.WaitGroup{}
	wg.Add(*ctx.Cmd.Workers)
	pool, worker, finish := cmd(ctx)
	jobs := make(chan Payload)

	for id := 1; id <= *ctx.Cmd.Workers; id++ {
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				message, info, err := worker(id, job)
				if err != nil {
					logger.JOB.Error(logger.Message{
						ShortMessage: fmt.Sprintf("Worker #%d, encountered and error!", id),
						FullMessage:  err.Error(),
						Params:       info,
					})
				} else {
					logger.JOB.Info(logger.Message{
						ShortMessage: fmt.Sprintf("Worker #%d, processed the unit!", id),
						FullMessage:  message,
						Params:       info,
					})
				}
			}
		}(id)
	}

	queue, err := pool()
	if err == nil {
		for len(queue) > 0 {
			for _, payload := range queue {
				jobs <- payload
			}

			queue, err = pool()
			if err != nil {
				logger.JOB.Error(logger.Message{
					ShortMessage: fmt.Sprintf("Job: '%s' exec stopped", *ctx.Cmd.Task),
					FullMessage:  err.Error(),
				})
				break
			}
		}
	} else {
		logger.JOB.Error(logger.Message{
			ShortMessage: fmt.Sprintf("Job: '%s' exec stopped", *ctx.Cmd.Task),
			FullMessage:  err.Error(),
		})
	}

	close(jobs)
	wg.Wait()
	ctx.State.Clear()

	if finish != nil {
		err = finish()
		if err != nil {
			logger.JOB.Error(logger.Message{
				ShortMessage: fmt.Sprintf("Job: '%s' exec stopped", *ctx.Cmd.Task),
				FullMessage:  err.Error(),
			})
		}
	}

	return nil
}
