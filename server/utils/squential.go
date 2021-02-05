// Package utils is used to store composition extensions for gRPC service servers only
package utils

import (
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Sequential extension for service to run method with certain params only once
type Sequential struct {
	runtime sync.Map
}

// Once don't allow calling certain method while previous call is running
// throw an error if it's called again while running
func (srv *Sequential) Once(key string, method func() error) error {
	if _, running := srv.runtime.Load(key); !running {
		srv.runtime.Store(key, true)
		defer srv.runtime.Delete(key)
		return method()
	}

	return status.Errorf(codes.ResourceExhausted, "job for '%s' in progress", key)
}
