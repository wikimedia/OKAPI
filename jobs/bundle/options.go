package bundle

import (
	"os"

	"okapi/helpers/writer"
	"okapi/lib/env"
	"okapi/lib/storage"
	"okapi/lib/task"
)

// Options bundle options struct
type Options struct {
	Limit    int
	Offset   int
	DumpName string
	DumpPath string
	Storage  storage.Connection
	Writer   *writer.Client
}

// Init function to setup initial data
func (opt *Options) Init(ctx *task.Context) error {
	opt.Limit, opt.Offset = *ctx.Cmd.Limit, *ctx.Cmd.Offset
	opt.Storage = storage.Local.Client()

	for _, path := range []string{"/exports", "/exports/" + ctx.Project.DBName} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(env.Context.VolumeMountPath+path, 0766)
		}
	}

	file, err := os.Create(ctx.Project.BundlePath())

	if err == nil {
		opt.Writer = writer.New(file)
	}

	return err
}
