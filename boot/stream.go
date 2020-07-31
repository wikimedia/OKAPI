package boot

import (
	"time"

	"okapi/lib/stream"
	"okapi/streams/page"
	"okapi/streams/revision"
)

// Stream function to start events server
func Stream() {
	clients := []*stream.Client{
		{
			Path:    "/revision-create",
			Handler: revision.Handler,
		},
		{
			Path:    "/page-delete",
			Handler: pageDelete.Handler,
		},
		{
			Path:    "/page-move",
			Handler: pageDelete.Handler,
		},
	}

	for _, client := range clients {
		go func(client *stream.Client) {
			for {
				stream.Subscribe(client)
				time.Sleep(2 * time.Second)
			}
		}(client)
	}

	select {}
}
