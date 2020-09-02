package bundle

import (
	"okapi/helpers/damaging"
	"okapi/helpers/projects"

	"okapi/helpers/writer"
	"okapi/lib/storage"
	"okapi/lib/task"
)

// Options bundle options struct
type Options struct {
	Limit      int
	Offset     int
	DumpName   string
	DumpPath   string
	RevDamaged []string
	Storage    storage.Connection
	Writer     *writer.Client
}

// Init function to setup initial data
func (opt *Options) Init(ctx *task.Context) error {
	opt.Limit, opt.Offset = *ctx.Cmd.Limit, *ctx.Cmd.Offset
	opt.Storage = storage.Local.Client()
	file, err := projects.CreateExportFile(ctx.Project)

	if err == nil {
		opt.Writer = writer.New(file)
	}

	opt.RevDamaged, err = damaging.Get(ctx.Project.DBName)

	if err == nil {
		err = damaging.Delete(ctx.Project.DBName)
	}

	return err
}
