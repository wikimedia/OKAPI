package storage

import (
	"okapi/lib/aws"
)

func remoteClient() Connection {
	return aws.NewStorage()
}
