package bundle

import (
	"okapi/helpers/bundle"
	"okapi/lib/task"
)

// Finish function to finish all the writing and send file to s3 bucket
func Finish(ctx *task.Context, options *Options) func() error {
	return func() error {
		options.Writer.Close()
		return bundle.Upload(ctx.Project)
	}
}
