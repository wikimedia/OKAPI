package bundle

import (
	"bufio"
	"os"
	"time"

	"okapi/helpers/progress"
	"okapi/lib/bzip"
	"okapi/lib/storage"
	"okapi/lib/task"
	"okapi/models"
)

// Finish function to finish all the writing and send file to s3 bucket
func Finish(ctx *task.Context, options *Options) func() error {
	return func() error {
		options.Writer.Close()

		path, err := bzip.Compress(options.DumpPath)
		if err != nil {
			return err
		}

		dumpFile, err := os.Open(path)
		defer dumpFile.Close()
		if err != nil {
			return err
		}

		fileInfo, err := dumpFile.Stat()
		if err != nil {
			return err
		}

		reader := &progress.Reader{
			Message: "File '" + options.DumpName + ".bz2' upload progress '%d%%'!",
			Reader:  bufio.NewReader(dumpFile),
			Size:    fileInfo.Size(),
		}

		if err = storage.Remote.Client().Put(options.DumpName+".bz2", reader); err != nil {
			return err
		}

		if ctx.Project.ID > 0 {
			ctx.Project.Size = ((float64)(reader.Size / 1024)) / 1024
			ctx.Project.Path = options.DumpName + ".bz2"
			ctx.Project.DumpedAt = time.Now()
			if err = models.Save(ctx.Project); err != nil {
				return err
			}
		}

		return nil
	}
}
