package storage

import (
	"io"
	"time"
)

// ConnectionName name of the storage connection
type ConnectionName string

// Connection interface to abstract data storing
type Connection interface {
	Get(path string) (io.ReadCloser, error)
	Put(path string, body io.Reader) error
	Link(path string, expire time.Duration) (string, error)
	Delete(path string) error
}
