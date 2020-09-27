package queue

import "okapi/lib/log"

// Context for queuer execution
type Context struct {
	Workers int
	Log     log.Log
}
