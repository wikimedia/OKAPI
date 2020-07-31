package runner

import "encoding/json"

// Command runner message
type Command struct {
	Task   string `json:"task"`
	DBName string `json:"db_name"`
}

// ParseCommand get command from json
func ParseCommand(info []byte) (*Command, error) {
	command := Command{}
	err := json.Unmarshal(info, &command)
	return &command, err
}
