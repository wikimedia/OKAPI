package processors

import (
	"okapi/lib/queue"
	"okapi/processors/page/delete"
	"okapi/processors/page/pull"
	"okapi/processors/page/revision"
)

// Workers list of workers for processors
var Workers = map[queue.Name]queue.Worker{
	queue.PageRevision: revision.Worker,
	queue.PagePull:     pull.Worker,
	queue.PageDelete:   delete.Worker,
}
