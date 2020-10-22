package pull

import (
	"okapi/lib/task"
	"okapi/models"
)

// Pool function to get new pages into the queue
func Pool(ctx *task.Context) func() ([]task.Payload, error) {
	limit, pointer := ctx.State.GetInt("limit", ctx.Params.Limit), ctx.State.GetInt("pointer", ctx.Params.Pointer)

	return func() ([]task.Payload, error) {
		pages := []*models.Page{}
		queue := []task.Payload{}
		query := models.DB().
			Model(&pages).
			Column("id", "project_id", "title", "site_url", "revision").
			Where("id > ?", pointer).
			Limit(limit)

		if ctx.Project.ID > 0 {
			query.Where("project_id = ?", ctx.Project.ID)
		}

		err := query.Order("id asc").Select()

		if err != nil {
			return queue, err
		}

		if len(pages) > 0 {
			ctx.State.Set("pointer", pointer)
			pointer = pages[len(pages)-1].ID
		}

		for _, page := range pages {
			page.Project = ctx.Project
			queue = append(queue, page)
		}

		return queue, nil
	}
}
