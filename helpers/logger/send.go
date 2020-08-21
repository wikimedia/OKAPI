package logger

import (
	"encoding/json"
	"sync"

	"okapi/lib/cmd"
	"okapi/lib/env"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client
var wg = sync.WaitGroup{}
var messages chan *Message = make(chan *Message)

// Client http client for graylog
func Client() *resty.Client {
	if client == nil {
		client = resty.New().
			SetHostURL(env.Context.LogURL).
			SetHeader("Content-Type", "application/json")
		go sender()
	}

	return client
}

// Close function to close sending channel
func Close() {
	wg.Wait()
	close(messages)
}

// Send message to graylog
func Send(message Message) {
	wg.Add(1)
	messages <- &message
}

// Send send log message to endpoint
func send(info map[string]interface{}) {
	message, _ := json.Marshal(info)
	Client().R().SetBody(message).Post("")
}

func sender() {
	for msg := range messages {
		req := map[string]interface{}{}

		if len(msg.Version) <= 0 {
			req["version"] = "1.1"
		}

		if len(msg.Host) <= 9 {
			req["host"] = env.Context.LogHost
		}

		if uint(msg.Level) == 0 {
			req["level"] = uint(INFO)
		} else {
			req["level"] = uint(msg.Level)
		}

		if *cmd.Context.Project != "*" {
			req["_project"] = *cmd.Context.Project
		}

		req["short_message"] = msg.ShortMessage
		req["full_message"] = msg.FullMessage
		req["_job"] = *cmd.Context.Task
		req["_server"] = *cmd.Context.Server
		req["_category"] = msg.Category

		if msg.Params != nil {
			for name, val := range msg.Params {
				req[name] = val
			}
		}

		wg.Done()
		send(req)
	}
}