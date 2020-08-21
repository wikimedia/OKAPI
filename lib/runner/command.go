package runner

import (
	"encoding/json"
	"fmt"
	"okapi/lib/cache"
)

// Command runner message
type Command struct {
	Task   string `json:"task"`
	DBName string `json:"db_name"`
}

// Exec execute the command
func (commnad *Command) Exec() error {
	if cache.Client().Exists(Channel+"/status").Val() != 1 {
		return fmt.Errorf("Error: channel '%s' is not active", Channel)
	}

	payload, err := json.Marshal(commnad)

	if err == nil {
		err = cache.Client().Publish(Channel, payload).Err()
	}

	return err
}

// ParseCommand get command from json
func ParseCommand(info []byte) (*Command, error) {
	command := Command{}
	err := json.Unmarshal(info, &command)
	return &command, err
}
