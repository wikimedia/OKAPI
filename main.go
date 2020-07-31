package main

import (
	"okapi/boot"
	"okapi/events"
	"okapi/helpers/logger"
	"okapi/lib/cache"
	"okapi/lib/cmd"
	"okapi/lib/env"
	"okapi/lib/minifier"
	"okapi/lib/storage"
	"okapi/listeners"
	"okapi/models"
)

func main() {
	cmd.Context.Parse()
	env.Context.Parse(".env")
	cache.Client()
	storage.Local.Client()
	storage.Remote.Client()
	minifier.Client()
	logger.Client()
	defer logger.Close()
	events.Init()
	listeners.Init()
	ch := cache.Client()
	db := models.DB()
	defer db.Close()
	defer ch.Close()

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
