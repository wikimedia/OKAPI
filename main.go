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
	"os"
)

func main() {
	os.Setenv("TZ", "UTC")
	cmd.Context.Parse()
	env.Context.Parse(".env")
	cache.Client()
	storage.Local.Client()
	storage.Remote.Client()
	minifier.Client()
	events.Init()
	listeners.Init()
	ch := cache.Client()
	db := models.DB()
	logger.Client()
	defer logger.Close()
	defer db.Close()
	defer ch.Close()

	// ORES
	if oresErr := ores.Client(); oresErr != nil {
		logger.API.Error(logger.Message{
			FullMessage: oresErr.Error(),
		})
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
