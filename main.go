package main

import (
	"okapi/boot"
	"okapi/events"
	"okapi/helpers/logger"
	"okapi/lib/cache"
	"okapi/lib/cmd"
	"okapi/lib/env"
	"okapi/lib/minifier"
	"okapi/lib/ores"
	"okapi/lib/storage"
	"okapi/listeners"
	"okapi/models"
	"okapi/validators"
	"os"
)

var libs = []func() error{
	cmd.Init,
	env.Init,
	cache.Init,
	storage.Init,
	minifier.Init,
	events.Init,
	listeners.Init,
	validators.Init,
	ores.Init,
	logger.Init,
}

func startup() error {
	os.Setenv("TZ", "UTC")

	for _, init := range libs {
		err := init()

		if err != nil {
			return err
		}
	}

	return nil
}

func cleanup() {
	models.DB().Close()
	cache.Client().Close()
	logger.Close()
}

func main() {
	err := startup()
	defer cleanup()

	if err != nil {
		logger.System.Panic("System startup failed", err.Error())
	}

	switch {
	case *cmd.Context.Task != "":
		boot.Task()
	case *cmd.Context.Server == string(cmd.Runner):
		boot.Runner()
	case *cmd.Context.Server == string(cmd.Stream):
		boot.Stream()
	case *cmd.Context.Server == string(cmd.Queue):
		boot.Queue()
	case *cmd.Context.Server == string(cmd.API):
		boot.API()
	}
}
