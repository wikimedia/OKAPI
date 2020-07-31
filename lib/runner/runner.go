package runner

import (
	"okapi/lib/cache"
)

// Executor struct to hold namespace and func of executable
type Executor struct {
	Namespace string
	Handler   func() error
}

// Run the process
func (exec *Executor) Run() {
	defer End.Send(exec.Namespace, &Message{
		Info: "Task execution finished!",
	})

	if cache.Client().Exists(exec.Namespace).Val() != 1 {
		defer cache.Client().Del(exec.Namespace).Result()

		Info.Send(exec.Namespace, &Message{
			Info: "Starting the task!",
		})

		err := exec.Handler()

		if err != nil {
			Error.Send(exec.Namespace, &Message{
				Info: err.Error(),
			})
		} else {
			Success.Send(exec.Namespace, &Message{
				Info: "Task successfully finished!",
			})
		}
	}
}
