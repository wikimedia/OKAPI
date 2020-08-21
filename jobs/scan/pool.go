package scan

import (
	"strings"

	"okapi/lib/dump"
	"okapi/lib/task"
	"okapi/models"
)

// Pool function to get new pages to the queue
func Pool(ctx *task.Context, options *Options, titles []string) func() ([]task.Payload, error) {
	titlesLength := len(titles)

	return func() ([]task.Payload, error) {
		queue := []task.Payload{}
		titlesSubset := []string{}

		titles, err := dump.Titles(ctx.Project.DBName, options.Folder)

		if err != nil {
			return queue, err
		}

		if options.Offset > titlesLength {
			return queue, nil
		}

		if options.Offset+options.Limit > titlesLength {
			titlesSubset = titles[options.Offset:titlesLength]
		} else {
			titlesSubset = titles[options.Offset : options.Offset+options.Limit]
		}

		for _, title := range titlesSubset {
			pageTitle := strings.Trim(title, " ")

			if len(pageTitle) > 0 {
				queue = append(queue, &models.Page{
					Title:     pageTitle,
					ProjectID: ctx.Project.ID,
					SiteURL:   ctx.Project.SiteURL,
					Project:   ctx.Project,
				})
			}
		}

		options.Offset += options.Limit

		return queue, err
	}
}
