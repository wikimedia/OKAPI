package run

import (
	"encoding/json"
	"fmt"
	"okapi/lib/cache"
)

// Cmd runner message
type Cmd struct {
	Task       string `json:"task"`
	DBName     string `json:"db_name"`
	Subscriber bool   `json:"subscriber"`
	Namespace  string `json:"-"`
}

// NewCmd create new command
func NewCmd(task string, dbName string, subscribe bool) *Cmd {
	return &Cmd{
		Task:       task,
		DBName:     dbName,
		Subscriber: subscribe,
		Namespace:  channel + "/" + task + "/" + dbName,
	}
}

// Enqueue add task to queue
func (cmd *Cmd) Enqueue() error {
	if !IsOnline() {
		return fmt.Errorf("channel '%s' is not active", channel)
	}

	payload, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	return cache.Client().Publish(channel, payload).Err()
}

// Success cmd was successfully executed
func (cmd *Cmd) Success() {
	Success.Send(cmd, &Msg{})
}

// Failed cmd failed to execute
func (cmd *Cmd) Failed(err error) {
	Failed.Send(cmd, &Msg{
		Info: err.Error(),
	})
}
