package boot

import (
	"time"

	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/stream"
	"okapi/streams"
)

// Stream function to start events server
func Stream() {
	ctx := stream.Context{
		Workers: *cmd.Context.Workers,
		Log:     logger.Steream,
	}

	for _, client := range streams.Clients {
		go func(client *stream.Client) {
			for {
				client.Subscribe(&ctx)
				time.Sleep(2 * time.Second)
			}
		}(client)
	}

	select {}
}
