package pagemove

import (
	"context"
	"log"
	"okapi-data-service/queues/pagedelete"
	"okapi-data-service/schema/v3"
	"okapi-data-service/streams/utils"
	"time"

	"github.com/go-redis/redis/v8"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
)

// Name segment of the stream in cache
const Name string = "stream/pagemove"

// Handler page move event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.PageMove) {
	return func(evt *eventstream.PageMove) {
		var err error

		if !utils.Exclude(evt.Data.Database) && utils.FilterNs(evt.Data.PageNamespace) {
			editor := &schema.Editor{
				Identifier: evt.Data.Performer.UserID,
				Name:       evt.Data.Performer.UserText,
				EditCount:  evt.Data.Performer.UserEditCount,
				Groups:     evt.Data.Performer.UserGroups,
				IsBot:      evt.Data.Performer.UserIsBot,
			}

			if !evt.Data.Performer.UserRegistrationDt.IsZero() {
				editor.DateStarted = &evt.Data.Performer.UserRegistrationDt
			}

			err = pagedelete.Enqueue(ctx, store, &pagedelete.Data{
				Title:  evt.Data.PriorState.PageTitle,
				DbName: evt.Data.Database,
				Editor: editor,
			})
		}

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}
