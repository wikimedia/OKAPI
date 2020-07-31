package scan

import (
	"strings"
	"sync"

	"okapi/helpers/logger"
	"okapi/lib/dump"
	"okapi/lib/task"
	"okapi/models"
)

// Pool function to get new pages to the queue
func Pool(ctx *task.Context, options *Options) func() ([]task.Payload, error) {
	return func() ([]task.Payload, error) {
		queue := []task.Payload{}
		projects := []models.Project{}
		query := models.DB().Model(&projects).Offset(options.Offset).Limit(options.Limit).Where("active = ?", true)

		if ctx.Project.ID > 0 {
			query.Where("id = ?", ctx.Project.ID)
		}

		if len(projects) > 0 {
			ctx.State.Set("limit", options.Limit)
			ctx.State.Set("offset", options.Offset)
		}

		options.Offset += options.Limit
		err := query.Order("id asc").Select()
		if err != nil {
			return queue, err
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(projects))
		for _, project := range projects {
			go func(project models.Project) {
				defer wg.Done()
				titles, err := dump.Titles(project.DBName, options.Folder)
				if err == nil {
					for _, title := range titles[options.Position:] {
						pageTitle := strings.Trim(title, " ")
						if len(pageTitle) > 0 {
							queue = append(queue, &models.Page{
								Title:     pageTitle,
								Lang:      project.Code,
								ProjectID: project.ID,
								SiteURL:   project.SiteURL,
							})
						}
					}
				} else {
					logger.JOB.Error(logger.Message{
						ShortMessage: "Job: 'bundle' exe failed in 'pool'",
						FullMessage:  err.Error(),
					})
				}
			}(project)
		}
		wg.Wait()

		return queue, err
	}
}
