package task

import (
	"okapi/lib/log"
	"okapi/models"
)

// Context task execution context
type Context struct {
	Project *models.Project
	Log     log.Log
	Params  Params
	State   State
}
