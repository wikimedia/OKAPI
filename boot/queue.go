package boot

import (
	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/queue"
	"okapi/processors"
)

// Queue function to start queue server
func Queue() {
	ctx := queue.Context{
		Workers: *cmd.Context.Workers,
		Log:     logger.Queue,
	}

	for queue, worker := range processors.Workers {
		go queue.Subscribe(&ctx, worker)
	}

	select {}
}
