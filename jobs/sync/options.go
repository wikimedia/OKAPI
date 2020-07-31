package sync

import (
	"okapi/lib/task"
)

// Options sync job options
type Options struct {
	Limit    int
	Offset   int
	Position int
}

// Init function to initialize options object
func (opt *Options) Init(ctx *task.Context) error {
	opt.Limit, opt.Offset = ctx.State.GetInt("limit", 100000), ctx.State.GetInt("offset", *ctx.Cmd.Position)
	opt.Position = opt.Offset
	return nil
}
