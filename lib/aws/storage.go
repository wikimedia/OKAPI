package aws

import (
	"io"
	"time"
)

// Storage implementation of storage interface
type Storage struct {
	Connection *Connection
}

// Get file from storage
func (store *Storage) Get(path string) (io.ReadCloser, error) {
	object, err := store.Connection.GetObjectFromS3(path)
	return object.Body, err
}

// Put put file into storage
func (store *Storage) Put(path string, body io.Reader) error {
	_, err := store.Connection.StreamFileToS3(path, body)
	return err
}

// Link get download link
func (store *Storage) Link(path string, expire time.Duration) (string, error) {
	req, _ := store.Connection.GetObjectRequestFromS3(path)
	return req.Presign(1 * time.Minute)
}

// Delete delete file from storage
func (store *Storage) Delete(path string) error {
	_, err := store.Connection.DeleteObjectFromS3(path)
	return err
}

// NewStorage creating news AWS storage
func NewStorage() *Storage {
	return &Storage{
		Connection: NewConnection(),
	}
}
