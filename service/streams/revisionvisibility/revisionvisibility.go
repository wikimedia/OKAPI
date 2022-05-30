package revisionvisibility

import (
	"context"
	"log"
	"okapi-data-service/queues/pagevisibility"
	"okapi-data-service/schema/v3"
	"okapi-data-service/streams/utils"
	"time"

	"github.com/go-redis/redis/v8"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
)

// Name segment of the stream in cache
const Name string = "stream/revisionvisibility"

// Handler revision visibility event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.RevisionVisibilityChange) {
	return func(evt *eventstream.RevisionVisibilityChange) {
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

			err = pagevisibility.Enqueue(ctx, store, &pagevisibility.Data{
				ID:         evt.Data.PageID,
				Title:      evt.Data.PageTitle,
				DbName:     evt.Data.Database,
				Revision:   evt.Data.RevID,
				RevisionDt: evt.Data.RevTimestamp,
				Visible:    evt.Data.Visibility.Text,
				Visibility: evt.Data.Visibility,
				Lang:       utils.Lang(evt.Data.Meta.Domain),
				SiteURL:    utils.SiteURL(evt.Data.Meta.Domain),
				Namespace:  evt.Data.PageNamespace,
				Editor:     editor,
			})
		}

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}
