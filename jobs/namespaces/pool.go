package namespaces

import (
	"okapi/lib/task"
	"okapi/models"
)

// Pool function to get projects into the queue
func Pool(ctx *task.Context) func() ([]task.Payload, error) {
	limit, offset := ctx.State.GetInt("limit", ctx.Params.Limit), ctx.State.GetInt("offset", ctx.Params.Offset)

	return func() ([]task.Payload, error) {
		queue := []task.Payload{}
		projects := []*models.Project{}

		query := models.DB().
			Model(&projects).
			Limit(limit).
			Offset(offset)

		offset += limit
		err := query.Order("id asc").Select()

		if err != nil {
			return queue, err
		}

		if len(projects) <= 0 {
			return queue, nil
		}

		for _, project := range projects {
			queue = append(queue, project)
		}

		return queue, nil
	}
}
