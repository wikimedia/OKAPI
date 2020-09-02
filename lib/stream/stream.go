package stream

import (
	"fmt"
	"sync"

	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/env"

	"github.com/r3labs/sse"
)

// Client struct for client info
type Client struct {
	Path    string
	Handler func(event *sse.Event)
}

// Subscribe method for subscribing to streams
func Subscribe(client *Client) {
	wg := sync.WaitGroup{}
	events := make(chan *sse.Event)
	httpClient := sse.NewClient(env.Context.StreamURL + client.Path)
	err := httpClient.SubscribeChan("messages", events)

	if err == nil {
		wg.Add(1)
		httpClient.OnDisconnect(func(c *sse.Client) {
			httpClient.Unsubscribe(events)
			wg.Done()
			logger.Steream.Error(fmt.Sprintf("Stream: '%s' stream was disconnected!", c.URL), "")
		})

		for i := 1; i <= *cmd.Context.Workers; i++ {
			go func() {
				for event := range events {
					client.Handler(event)
				}
			}()
		}
	} else {
		logger.Steream.Error(fmt.Sprintf("Stream '%s' connection failed", client.Path), err.Error())
	}

	wg.Wait()
}
