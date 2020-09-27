package run

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// NewRunner create new runner with given cache client
func NewRunner(sub *redis.Client) *Runner {
	return &Runner{
		sub: sub,
	}
}

// Runner main execution unit
type Runner struct {
	sub *redis.Client
}

// NewCmd create cmd from redis message
func (run *Runner) NewCmd(msg *redis.Message) (*Cmd, error) {
	cmd := &Cmd{}
	err := json.Unmarshal([]byte(msg.Payload), cmd)

	if err == nil {
		cmd.Namespace = channel + "/" + cmd.Task + "/" + cmd.DBName
	}

	return cmd, err
}

// Channel get messages channel
func (run *Runner) Channel() <-chan *redis.Message {
	return run.sub.Subscribe(channel).Channel()
}

// Connect connect to subscriber
func (run *Runner) Connect(cmd *Cmd) error {
	if cmd.Subscriber {
		retry := retries

		for {
			time.Sleep(1 * time.Second)
			channels, err := run.sub.PubSubChannels(cmd.Namespace).Result()

			if err != nil || len(channels) > 0 || retry <= 0 {
				return err
			}

			retry--
		}
	}

	return nil
}
