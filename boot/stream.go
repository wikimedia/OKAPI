package boot

import (
	"time"

	"okapi/lib/stream"
	"okapi/streams"
)

// Stream function to start events server
func Stream() {
	for _, client := range streams.Clients {
		go func(client *stream.Client) {
			for {
				stream.Subscribe(client)
				time.Sleep(2 * time.Second)
			}
		}(client)
	}

	select {}
}
