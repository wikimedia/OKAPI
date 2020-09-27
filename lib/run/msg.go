package run

import (
	"encoding/json"
	"okapi/lib/cache"
)

// Msg message for pub/sub communication
type Msg struct {
	Status Status      `json:"status"`
	Info   interface{} `json:"info"`
}

// Send message into namespace
func (msg *Msg) Send(namespace string) error {
	info, err := json.Marshal(msg)

	if err == nil {
		err = cache.Client().Publish(namespace, info).Err()
	}

	return err
}
