package task

import (
	"fmt"
	"sync"
)

// Name task name
type Name string

// Task struct
type Task func(ctx *Context) (Pool, Worker, Finish, error)

// Exec some task with context
func Exec(cmd Task, ctx *Context) error {
	if ctx.Params.Restart {
		ctx.State.Clear()
	}
	defer ctx.State.Clear()

	pool, worker, finish, err := cmd(ctx)

	if err != nil {
		return err
	}

	if pool != nil && worker != nil {
		wg := &sync.WaitGroup{}
		wg.Add(ctx.Params.Workers)
		jobs := make(chan Payload)

		for id := 1; id <= ctx.Params.Workers; id++ {
			go runWorker(ctx, id, worker, wg, jobs)
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

func runWorker(ctx *Context, id int, handler Worker, wg *sync.WaitGroup, jobs chan Payload) {
	defer wg.Done()

	for job := range jobs {
		message, info, err := handler(id, job)

		if err != nil {
			ctx.Log.Error(fmt.Sprintf("worker #%d error db_name: '%s'", id, ctx.Params.DBName), err.Error(), info)
		} else {
			ctx.Log.Info(fmt.Sprintf("worker #%d processed db_name: '%s'", id, ctx.Params.DBName), message)
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
