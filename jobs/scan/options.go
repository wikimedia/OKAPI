package scan

import (
	"fmt"
	"okapi/lib/dump"
	"okapi/lib/task"
)

// Options scan job options
type Options struct {
	Offset   int
	Limit    int
	Position int
	Folder   string
}

// Init function to get initial data for the job
func (opt *Options) Init(ctx *task.Context) error {
	var err error
	opt.Folder = ctx.State.GetString("folder", "")
	opt.Limit, opt.Offset = ctx.State.GetInt("limit", ctx.Params.Limit), ctx.State.GetInt("offset", ctx.Params.Offset)
	opt.Position = opt.Offset

	if ctx.Project.ID <= 0 {
		return fmt.Errorf("this task requires project to run")
	}

	if len(opt.Folder) > 0 {
		return nil
	}

	opt.Folder, err = dump.Folder()

	if err == nil {
		ctx.State.Set("folder", opt.Folder)
	}

	return err
}
