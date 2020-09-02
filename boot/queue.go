package boot

import (
	"okapi/lib/cmd"
	"okapi/lib/queue"
	"okapi/processors"
)

// Queue function to start queue server
func Queue() {
	for subscriber, worker := range processors.Workers {
		go queue.Subscribe(subscriber, worker, *cmd.Context.Workers)
	}

	select {}
}
