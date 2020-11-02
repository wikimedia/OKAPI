package export

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"okapi/lib/task"
	"os"
	"sync"
)

func readWorker(ctx *task.Context, paths chan *path, files chan *file, wg *sync.WaitGroup) {
	for path := range paths {
		buffer, err := readHTML(path.full)

		if err != nil {
			ctx.Log.Error(fmt.Sprintf("error reading file '%s'", path), err.Error())
			continue
		}

		wg.Add(1)
		files <- &file{
			name:   path.file,
			buffer: buffer,
		}

		ctx.Log.Info("added to write queue", fmt.Sprintf("file: '%s'", path))
	}

	wg.Done()
}

func readHTML(path string) (*bytes.Buffer, error) {
	ref, err := os.Open(path)
	defer ref.Close()

	if err != nil {
		return nil, err
	}

	html, err := ioutil.ReadAll(ref)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(html), nil
}
