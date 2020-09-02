package runner

import (
	"encoding/json"
	"fmt"
	"okapi/jobs"
	"okapi/lib/cache"
	"okapi/lib/task"
	"time"

	"github.com/go-redis/redis"
)

const connectionRetries = 10

// Command runner message
type Command struct {
	Task       string `json:"task"`
	DBName     string `json:"db_name"`
	Subscriber bool   `json:"subscriber"`
	Namespace  string `json:"-"`
}

// Error notify subscriber about error
func (cmd *Command) Error(message string) {
	Error.Send(cmd.Namespace, &Message{
		Info: message,
	})
}

// Info send info to subscriber
func (cmd *Command) Info(message string) {
	Info.Send(cmd.Namespace, &Message{
		Info: message,
	})
}

// Success send success message
func (cmd *Command) Success(message string) {
	Success.Send(cmd.Namespace, &Message{
		Info: message,
	})
}

// End send end message to subscriber for him to unsubscribe
func (cmd *Command) End(message string) {
	End.Send(cmd.Namespace, &Message{
		Info: message,
	})
}

// Job get tasks executor
func (cmd *Command) Job() (task.Task, error) {
	job, exists := jobs.Tasks[task.Name(cmd.Task)]

	if !exists {
		return nil, fmt.Errorf("task '%s' not found", cmd.Task)
	}

	return job, nil
}

// Stop clear namespace for other commands of the same type
func (cmd *Command) Stop() {
	cache.Client().Del(cmd.Namespace)
}

// Running set cache namespace to identify that command is running
func (cmd *Command) Running() {
	cache.Client().Set(cmd.Namespace, "", 24*time.Hour)
}

// IsRunning check wether this command already running
func (cmd *Command) IsRunning() bool {
	return cache.Client().Exists(cmd.Namespace).Val() == 1
}

// WaitSubscriber wait for subscriber to connect to ensure he's listening to the communication
func (cmd *Command) WaitSubscriber(sub *redis.Client) error {
	if cmd.Subscriber {
		retry := connectionRetries

		for {
			time.Sleep(1 * time.Second)
			channels, err := sub.PubSubChannels(cmd.Namespace).Result()

			if err != nil || len(channels) > 0 || retry <= 0 {
				return err
			}

			retry--
		}
	}

	return nil
}

// Exec execute the command
func (cmd *Command) Exec() error {
	if !IsOnline() {
		return fmt.Errorf("Channel '%s' is not active", Channel)
	}

	payload, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	return cache.Client().Publish(Channel, payload).Err()
}

// ParseCommand get command from json
func ParseCommand(info []byte) (*Command, error) {
	command := Command{}
	err := json.Unmarshal(info, &command)

	if err == nil {
		command.Namespace = Channel + "/" + command.Task + "/" + command.DBName
	}

	return &command, err
}
