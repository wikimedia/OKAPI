package stream

import (
	"fmt"
	"okapi/lib/env"
	"sync"

	"github.com/r3labs/sse"
)

// Client struct for client info
type Client struct {
	Path    string
	Handler func(event *sse.Event)
}

// Subscribe from events stream
func (client *Client) Subscribe(ctx *Context) {
	url := env.Context.StreamURL + client.Path
	http := sse.NewClient(url)
	events := make(chan *sse.Event)
	wg := sync.WaitGroup{}

	err := http.SubscribeChan("messages", events)

	if err == nil {
		wg.Add(1)

		http.OnDisconnect(func(c *sse.Client) {
			http.Unsubscribe(events)
			wg.Done()
			ctx.Log.Error("stream disconnected", fmt.Sprintf("url: '%s'", url))
		})

		for i := 1; i <= ctx.Workers; i++ {
			go func() {
				for event := range events {
					client.Handler(event)
				}
			}()
		}
	} else {
		ctx.Log.Error("stream failed to connect", fmt.Sprintf("url: '%s'", url))
	}

	wg.Wait()
}
