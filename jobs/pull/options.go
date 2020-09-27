package pull

import (
	"fmt"
	"okapi/lib/task"
)

// Options pull job options
type Options struct {
	Limit    int
	Offset   int
	Position int
}

// Init function to initialize options object
func (opt *Options) Init(ctx *task.Context) error {
	opt.Limit, opt.Offset = ctx.State.GetInt("limit", ctx.Params.Limit), ctx.State.GetInt("offset", ctx.Params.Offset)
	opt.Position = opt.Offset

	if ctx.Project.ID <= 0 {
		return fmt.Errorf("this task requires project to run")
	}

	return nil
}
