package main

import (
	"okapi/boot"
	"okapi/events"
	"okapi/helpers/logger"
	"okapi/lib/cache"
	"okapi/lib/cmd"
	"okapi/lib/elastic"
	"okapi/lib/env"
	"okapi/lib/minifier"
	"okapi/lib/ores"
	"okapi/lib/storage"
	"okapi/listeners"
	"okapi/models"
	"okapi/validators"
	"os"
)

var servers = map[cmd.Server]func(){
	cmd.Runner: boot.Runner,
	cmd.Stream: boot.Stream,
	cmd.Queue:  boot.Queue,
	cmd.API:    boot.API,
}

var libs = []func() error{
	cmd.Init,
	env.Init,
	models.Init,
	cache.Init,
	storage.Init,
	minifier.Init,
	events.Init,
	listeners.Init,
	validators.Init,
	ores.Init,
	logger.Init,
	elastic.Init,
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
	models.Close()
	cache.Close()
	elastic.Close()
	logger.Close()
}

func main() {
	err := startup()
	defer cleanup()

	if err != nil {
		logger.System.Panic("System startup failed", err.Error())
	}

	if len(*cmd.Context.Task) > 0 {
		boot.Task()
	} else if serve, ok := servers[cmd.Server(*cmd.Context.Server)]; ok {
		serve()
	}
}
