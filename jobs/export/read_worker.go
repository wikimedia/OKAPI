package export

import (
	"bytes"
	"io/ioutil"
	"okapi/helpers/logger"
	"os"
	"sync"
)

func readWorker(paths chan string, length int, files chan *file, group *sync.WaitGroup) {
	group.Add(1)

	for path := range paths {
		ref, err := os.Open(path)
		defer ref.Close()

		if err != nil {
			logger.Job.Log(err, "Error opening file")
			continue
		}

		html, err := ioutil.ReadAll(ref)

		if err != nil {
			logger.Job.Log(err, "Error reading the file")
			continue
		}

		group.Add(1)
		files <- &file{
			name:   path[length:],
			buffer: bytes.NewBuffer(html),
		}
	}

	group.Done()
}
