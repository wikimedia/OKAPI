package index

import (
	"math"
	"okapi/lib/task"
	"okapi/models"
)

// Pool function for queriing pages from database inside the task
func Pool(ctx *task.Context) func() ([]task.Payload, error) {
	limit, offset := ctx.State.GetInt("limit", ctx.Params.Limit), ctx.State.GetInt("offset", ctx.Params.Offset)

	return func() ([]task.Payload, error) {
		queue := []task.Payload{}
		pages := []*models.Page{}

		query := models.DB().
			Model(&pages).
			Limit(limit).
			Offset(offset)

		if ctx.Project.ID > 0 {
			query.Where("project_id = ?", ctx.Project.ID)
		}

		offset += limit
		err := query.Order("id asc").Select()

		if err != nil {
			return queue, err
		}

		if len(pages) <= 0 {
			return queue, nil
		}

		length := int(math.Ceil(float64(len(pages)) / float64(ctx.Params.Workers)))

		for i := 0; i < ctx.Params.Workers; i++ {
			start, end := i*length, ((i * length) + length)

			if end > len(pages) {
				end = len(pages)
			}

			chunk := pages[start:end]

			if len(chunk) > 0 {
				queue = append(queue, chunk)
			}
		}

		return queue, nil
	}
}
