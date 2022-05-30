package page

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

type store interface {
	storage.Putter
	storage.Deleter
	storage.Getter
}

type Storage struct {
	Local  store
	Remote store
}

// Get get data from the storage
func (s Storage) Get(path string) (io.ReadCloser, error) {
	return s.Local.Get(path)
}

// Put upload file to remote and local storages
func (s Storage) Put(path string, body io.Reader) error {
	errs := make(chan error, 2)
	data, err := ioutil.ReadAll(body) // TODO change params from io.Reader to []bytes

	if err != nil {
		return err
	}

	go func() {
		errs <- s.Local.Put(path, bytes.NewReader(data))
	}()

	go func() {
		errs <- s.Remote.Put(fmt.Sprintf("page/%s", path), bytes.NewReader(data))
	}()

	for i := 0; i < 2; i++ {
		err := <-errs

		if err != nil {
			return err
		}
	}

	return nil
}

// Delete remove file to remote and local storages
func (s Storage) Delete(path string) error {
	errs := make(chan error, 2)

	go func() {
		errs <- s.Local.Delete(path)
	}()

	go func() {
		errs <- s.Remote.Delete(fmt.Sprintf("page/%s", path))
	}()

	for i := 0; i < 2; i++ {
		err := <-errs

		if err != nil {
			return err
		}
	}

	return nil
}
