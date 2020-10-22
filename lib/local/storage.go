package local

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"okapi/lib/env"
)

// Storage implementation of storage interface
type Storage struct{}

func getPath(path string) string {
	return env.Context.VolumeMountPath + "/" + regexp.MustCompile(`//`).ReplaceAllString(path, "/")
}

// Get file from storage
func (store *Storage) Get(path string) (io.ReadCloser, error) {
	return os.Open(getPath(path))
}

// Put put file into storage
func (store *Storage) Put(path string, body io.Reader) error {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	localPath := getPath(path)
	dir, _ := filepath.Split(localPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
	}

	if err != nil {
		return err
	}

	return ioutil.WriteFile(localPath, buf, 0766)
}

// Link get download link
func (store *Storage) Link(path string, expire time.Duration) (string, error) {
	return getPath(path), nil
}

// Delete delete file from storage
func (store *Storage) Delete(path string) error {
	if len(path) <= 0 {
		return nil
	}

	return os.Remove(getPath(path))
}

// NewStorage creating news local storage
func NewStorage() *Storage {
	return &Storage{}
}
