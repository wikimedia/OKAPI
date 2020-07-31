package storage

import (
	"okapi/lib/local"
)

func localClient() Connection {
	return local.NewStorage()
}
