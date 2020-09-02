package boot

import (
	"okapi/helpers/logger"
	"okapi/lib/cache"
	"okapi/lib/runner"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
)

// Runner function to run tasks through redis
func Runner() {
	sub := cache.NewClient()
	defer sub.Close()
	defer runner.Offline()
	runner.Online()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go interrupt(signals)

	channel := sub.Subscribe(runner.Channel).Channel()

	for message := range channel {
		cmd, err := runner.ParseCommand([]byte(message.Payload))

		if err != nil {
			logger.Runner.Error("Command parse failed", err.Error())
			continue
		}

		go execute(cmd, sub)
	}
}

func execute(cmd *runner.Command, sub *redis.Client) {
	err := cmd.WaitSubscriber(sub)

	if err != nil {
		logger.Runner.Error("Subscriber wait failed", err.Error())
		return
	}

	err = runner.Execute(cmd)

	if err != nil {
		cmd.Error(err.Error())
		logger.Runner.Error("Runner exec failed", err.Error())
	} else {
		cmd.Success("Runner finished the task")
	}
}

func interrupt(signals chan os.Signal) {
	for range signals {
		runner.Offline()
		os.Exit(0)
	}
}
