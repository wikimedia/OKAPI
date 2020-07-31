package scan

import (
	"okapi/lib/dump"
	"okapi/lib/task"
)

// Options scan job options
type Options struct {
	Limit    int
	Offset   int
	Position int
	Folder   string
}

// Init function to get initial data for the job
func (opt *Options) Init(ctx *task.Context) error {
	opt.Limit, opt.Offset = ctx.State.GetInt("limit", 1), ctx.State.GetInt("offset", 0)
	opt.Folder = ctx.State.GetString("folder", "")
	opt.Position = ctx.State.GetInt("position", *ctx.Cmd.Position)

	if len(opt.Folder) > 0 {
		return nil
	}

	var err error
	opt.Folder, err = dump.Folder()
	if err == nil {
		ctx.State.Set("folder", opt.Folder)
	}

	return err
}
