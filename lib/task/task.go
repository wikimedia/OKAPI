package task

import (
	"fmt"
	"sync"

	"okapi/helpers/logger"
	"okapi/helpers/state"
	"okapi/lib/cmd"
	"okapi/models"

	"gopkg.in/gookit/color.v1"
)

// Context task execution context
type Context struct {
	Cmd     *cmd.Params
	State   *state.Client
	Project *models.Project
}

// Task struct
type Task func(ctx *Context) (Pool, Worker, Finish, error)

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
	defer func() { ctx.State.Clear() }()
	if *ctx.Cmd.Restart {
		ctx.State.Clear()
	}

	pool, worker, finish, err := cmd(ctx)

	if err != nil {
		return err
	}

	if pool != nil && worker != nil {
		wg := &sync.WaitGroup{}
		wg.Add(*ctx.Cmd.Workers)
		jobs := make(chan Payload)

		for id := 1; id <= *ctx.Cmd.Workers; id++ {
			go runWorker(id, worker, wg, jobs)
		}

		err := runPool(pool, jobs)
		close(jobs)
		wg.Wait()

		if err != nil {
			return err
		}
	}

	if finish != nil {
		err := finish()

		if err != nil {
			return err
		}
	}

	return nil
}

func runWorker(id int, handler Worker, wg *sync.WaitGroup, jobs chan Payload) {
	defer wg.Done()

	for job := range jobs {
		message, info, err := handler(id, job)

		if err != nil {
			logger.Job.Error(fmt.Sprintf("Worker #%d - encountered and error!", id), err.Error(), info)
		} else {
			color.Success.Println(fmt.Sprintf("Worker #%d - %s", id, message))
		}
	}
}

func runPool(pool Pool, jobs chan Payload) error {
	wg := sync.WaitGroup{}
	preload := 0
	queue, err := pool()

	if err != nil {
		return err
	}

	for len(queue) > 0 {
		if preload > 5 {
			wg.Wait()
		}

		wg.Add(1)
		preload++

		go func(queue []Payload) {
			for _, payload := range queue {
				jobs <- payload
			}

			wg.Done()
			preload--
		}(queue)

		queue, err = pool()

		if err != nil {
			return err
		}
	}

	wg.Wait()

	return nil
}
