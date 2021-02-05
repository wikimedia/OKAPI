package pagedelete

import (
	"context"
	"log"
	"okapi-data-service/queues/pagedelete"
	"okapi-data-service/streams/utils"
	"time"

	"github.com/go-redis/redis/v8"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
)

// Name segment of the stream in cache
const Name string = "stream/pagedelete"

// Handler page delete event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.PageDelete) {
	return func(evt *eventstream.PageDelete) {
		var err error

		if !utils.Exclude(evt.Data.Database) {
			err = pagedelete.Enqueue(ctx, store, &pagedelete.Data{
				Title:  evt.Data.PageTitle,
				DbName: evt.Data.Database,
			})
		}

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}