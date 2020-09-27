package boot

import (
	"okapi/helpers/jobs"
	"okapi/helpers/logger"
	"okapi/lib/cache"
	lib_cmd "okapi/lib/cmd"
	"okapi/lib/run"
	"okapi/lib/task"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
)

// Runner function to run tasks through pub/sub
func Runner() {
	sub := cache.NewClient()
	runner := run.NewRunner(sub)

	run.SetOnline()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)
	go interrupt(signals, sub)

	for msg := range runner.Channel() {
		cmd, err := runner.NewCmd(msg)

		if err != nil {
			logger.Runner.Error("cmd parse failed", err.Error())
			continue
		}

		err = runner.Connect(cmd)

		if err != nil {
			logger.Runner.Error("can't connect to remote", err.Error())
			continue
		}

		job, ctx, err := jobs.FromCMD(cmd, lib_cmd.Context)

		if err != nil {
			logger.Runner.Error("job init failed", err.Error())
			continue
		}

		go execute(job, &ctx, cmd)
	}
}

func execute(job task.Task, ctx *task.Context, cmd *run.Cmd) {
	defer cmd.Success()
	err := task.Exec(job, ctx)

	if err != nil {
		logger.Runner.Error("job exec failed", err.Error())
		cmd.Failed(err)
	}
}

func interrupt(signals chan os.Signal, sub *redis.Client) {
	for range signals {
		run.SetOffline()
		sub.Close()
		os.Exit(0)
	}
}
