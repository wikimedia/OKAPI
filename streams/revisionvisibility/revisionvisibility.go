package revisionvisibility

import (
	"context"
	"log"
	"okapi-data-service/queues/pagevisibility"
	"time"

	"github.com/go-redis/redis/v8"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
)

// Name segment of the stream in cache
const Name string = "stream/revisionvisibility"

// Handler revision visibility event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.RevisionVisibilityChange) {
	return func(evt *eventstream.RevisionVisibilityChange) {
		err := pagevisibility.Enqueue(ctx, store, &pagevisibility.Data{
			Title:    evt.Data.PageTitle,
			DbName:   evt.Data.Database,
			Revision: evt.Data.RevID,
		})

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}
