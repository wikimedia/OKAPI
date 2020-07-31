package runner

import (
	"encoding/json"
	"okapi/lib/cache"
)

// Status run status
type Status string

// ToString convert message to string with current status
func (status Status) ToString(message *Message) (string, error) {
	message.Status = status
	return message.ToString()
}

// Send message with certain status
func (status Status) Send(namespace string, Message *Message) error {
	Message.Status = status
	return Message.Send(namespace)
}

// Statuses for runner
const (
	Info    Status = "info"
	Error   Status = "error"
	Success Status = "success"
	End     Status = "end"
)

// Message status message
type Message struct {
	Status Status `json:"status"`
	Info   string `json:"info"`
}

// ToString convert message to string
func (message *Message) ToString() (string, error) {
	info, err := json.Marshal(message)
	return string(info), err
}

// Send message into namespace
func (message *Message) Send(namespace string) error {
	payload, err := message.ToString()

	if err == nil {
		err = cache.Client().Publish(namespace, payload).Err()
	}

	return err
}
