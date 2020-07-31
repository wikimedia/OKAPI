package boot

import (
	"okapi/lib/queue"
	"okapi/processors/page/delete"
	"okapi/processors/scan"
	"okapi/processors/sync"
)

// Queue function to start queue server
func Queue() {
	queues := map[queue.Name]queue.Worker{
		queue.Sync:       sync.Worker,
		queue.Scan:       scan.Worker,
		queue.DeletePage: deletePage.Worker,
	}

	for subscriber, worker := range queues {
		go queue.Subscribe(subscriber, worker)
	}

	select {}
}
