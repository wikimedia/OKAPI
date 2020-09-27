package export

import (
	"archive/tar"
	"fmt"
	"io"
	"okapi/lib/task"
	"sync"
)

func writeWorker(ctx *task.Context, files chan *file, tr *tar.Writer, wg *sync.WaitGroup) {
	for file := range files {
		err := tr.WriteHeader(&tar.Header{
			Name: file.name,
			Size: int64(file.buffer.Len()),
			Mode: 0766,
		})

		if err != nil {
			ctx.Log.Error("write failed", err.Error())
		}

		_, err = io.Copy(tr, file.buffer)

		if err != nil {
			ctx.Log.Error("write failed", err.Error())
		}

		wg.Done()
		ctx.Log.Info("written a file", fmt.Sprintf("file: '%s'", file.name))
	}
}
