package streams

import (
	"okapi/lib/stream"
	"okapi/streams/page/delete"
	"okapi/streams/page/revision"
)

// Clients stream clients handler
var Clients = []*stream.Client{
	{
		Path:    "/revision-create",
		Handler: revision.Handler,
	},
	{
		Path:    "/page-delete",
		Handler: delete.Handler,
	},
	{
		Path:    "/page-move",
		Handler: delete.Handler,
	},
}
