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
	opt.Limit, opt.Offset = 100000, 0
	opt.DumpName = "dump_" + *ctx.Cmd.Project + ".tar"
	opt.DumpPath = env.Context.VolumeMountPath + "/" + opt.DumpName
	opt.Storage = storage.Local.Client()

	file, err := os.Create(opt.DumpPath)
	if err != nil {
		return err
	}

	opt.Writer = writer.New(file)
	return nil
}
