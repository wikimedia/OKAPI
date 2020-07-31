package stream

import (
	"fmt"
	"sync"

	"github.com/r3labs/sse"
	"okapi/helpers/logger"
	"okapi/lib/cmd"
	"okapi/lib/env"
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
			logger.STREAM.Error(logger.Message{
				ShortMessage: "Stream: 'messages' stream was disconnected!",
				FullMessage:  c.URL,
			})
		})

		for i := 1; i <= *cmd.Context.Workers; i++ {
			go func() {
				for event := range events {
					client.Handler(event)
				}
			}()
		}
	} else {
		logger.STREAM.Error(logger.Message{
			ShortMessage: fmt.Sprintf("Stream: '%s' connection failed", client.Path),
			FullMessage:  err.Error(),
		})
	}

	wg.Wait()
}
