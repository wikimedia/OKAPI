package revisioncreate

import (
	"context"
	"log"
	"okapi-data-service/queues/pagepull"
	"okapi-data-service/streams/utils"
	"time"

	"github.com/go-redis/redis/v8"
	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	ores "github.com/protsack-stephan/mediawiki-ores-client"
)

// Name segment of the stream in cache
const Name string = "stream/revisioncreate"

// Handler revision create event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.RevisionCreate) {
	return func(evt *eventstream.RevisionCreate) {
		var err error

		if !evt.Data.PageIsRedirect && !ores.ModelDamaging.Supports(evt.Data.Database) && !utils.Exclude(evt.Data.Database) && utils.FilterNs(evt.Data.PageNamespace) {
			err = pagepull.Enqueue(ctx, store, &pagepull.Data{
				Title:   evt.Data.PageTitle,
				DbName:  evt.Data.Database,
				Lang:    utils.Lang(evt.Data.Meta.Domain),
				SiteURL: utils.SiteURL(evt.Data.Meta.Domain),
			})
		}

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}
