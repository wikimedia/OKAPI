package bundle

import (
	"okapi/lib/task"
	"okapi/models"

	"github.com/go-pg/pg/v9"
)

// Pool adding new pages to the queue
func Pool(ctx *task.Context, options *Options) func() ([]task.Payload, error) {
	return func() ([]task.Payload, error) {
		pages := []*models.Page{}
		queue := []task.Payload{}
		query := models.DB().Model(&pages).
			Offset(options.Offset).
			Limit(options.Limit).
			Where("path is not null")

		if ctx.Project.ID > 0 {
			query.Where("project_id = ?", ctx.Project.ID)
		}

		if len(options.RevDamaged) > 0 {
			query.Where("revision not in (?)", pg.In(options.RevDamaged))
		}

		options.Offset += options.Limit
		err := query.Order("id asc").Select()

		if err != nil {
			return queue, err
		}

		for _, page := range pages {
			queue = append(queue, page)
		}

		return queue, nil
	}
}
