package export

import (
	"archive/tar"
	"io"
	"okapi/helpers/logger"
	"sync"
)

func writeWorker(files chan *file, tr *tar.Writer, wg *sync.WaitGroup) {
	for file := range files {
		err := tr.WriteHeader(&tar.Header{
			Name: file.name,
			Size: int64(file.buffer.Len()),
			Mode: 0766,
		})

		if err != nil {
			logger.Job.Error("Worker write failed", err.Error())
		}

		_, err = io.Copy(tr, file.buffer)

		if err != nil {
			logger.Job.Error("Write worker failed", err.Error())
		}

		wg.Done()
	}
}
