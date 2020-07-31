package sync

import (
	"okapi/lib/task"
	"okapi/models"
)

// Pool function to get new pages into the queue
func Pool(ctx *task.Context, options *Options) func() ([]task.Payload, error) {
	return func() ([]task.Payload, error) {
		pages := []*models.Page{}
		queue := []task.Payload{}
		query := models.DB().Model(&pages).Offset(options.Offset).Limit(options.Limit)

		if ctx.Project.ID > 0 {
			query.Where("project_id = ?", ctx.Project.ID)
		}

		if len(pages) > 0 {
			ctx.State.Set("limit", options.Limit)
			ctx.State.Set("offset", options.Offset)
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
