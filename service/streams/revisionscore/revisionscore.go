package revisionscore

import (
	"context"
	"log"
	"okapi-data-service/queues/pagefetch"
	"okapi-data-service/schema/v3"
	"okapi-data-service/streams/utils"
	"time"

	"github.com/go-redis/redis/v8"

	eventstream "github.com/protsack-stephan/mediawiki-eventstream-client"
	ores "github.com/protsack-stephan/mediawiki-ores-client"
)

// Name segment of the stream in cache
const Name string = "stream/revisionscore"

// Handler revision score event handler
func Handler(ctx context.Context, store redis.Cmdable, expire time.Duration) func(evt *eventstream.RevisionScore) {
	return func(evt *eventstream.RevisionScore) {
		var err error

		if !evt.Data.PageIsRedirect && ores.ModelDamaging.Supports(evt.Data.Database) && !utils.Exclude(evt.Data.Database) && utils.FilterNs(evt.Data.PageNamespace) {
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

			data := &pagefetch.Data{
				Title:     evt.Data.PageTitle,
				DbName:    evt.Data.Database,
				Lang:      utils.Lang(evt.Data.Meta.Domain),
				SiteURL:   utils.SiteURL(evt.Data.Meta.Domain),
				Revision:  evt.Data.RevID,
				Namespace: evt.Data.PageNamespace,
				Scores:    &schema.Scores{},
				Editor:    editor,
			}

			if evt.Data.Scores.Damaging.Probability.True != 0 || evt.Data.Scores.Damaging.Probability.False != 0 {
				data.Scores.Damaging = &ores.ScoreDamaging{}
				data.Scores.Damaging.Prediction = evt.Data.Scores.Damaging.Probability.True > evt.Data.Scores.Damaging.Probability.False
				data.Scores.Damaging.Probability.True = evt.Data.Scores.Damaging.Probability.True
				data.Scores.Damaging.Probability.False = evt.Data.Scores.Damaging.Probability.False
			}

			if evt.Data.Scores.Goodfaith.Probability.True != 0 || evt.Data.Scores.Goodfaith.Probability.False != 0 {
				data.Scores.GoodFaith = &ores.ScoreGoodFaith{}
				data.Scores.GoodFaith.Prediction = evt.Data.Scores.Goodfaith.Probability.True > evt.Data.Scores.Goodfaith.Probability.False
				data.Scores.GoodFaith.Probability.True = evt.Data.Scores.Goodfaith.Probability.True
				data.Scores.GoodFaith.Probability.False = evt.Data.Scores.Goodfaith.Probability.False
			}

			err = pagefetch.Enqueue(ctx, store, data)
		}

		if err != nil {
			log.Printf("%s: %v\n", Name, err)
		} else if err := store.Set(ctx, Name, evt.Data.Meta.Dt, expire).Err(); err != nil {
			log.Printf("%s: %v\n", Name, err)
		}
	}
}
