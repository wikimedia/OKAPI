package projects

import (
	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Pool function to get projects into the queue
func Pool(body *wiki.Projects) func() ([]task.Payload, error) {
	return func() ([]task.Payload, error) {
		queue := []task.Payload{}

		if body.Sitematrix != nil {
			for _, project := range body.Sitematrix {
				for _, site := range project.Site {
					queue = append(queue, &models.Project{
						Name:      project.Name,
						Code:      project.Code,
						SiteName:  site.SiteName,
						SiteURL:   site.URL,
						SiteCode:  site.Code,
						DBName:    site.DBName,
						Dir:       project.Dir,
						LocalName: project.LocalName,
						Active:    !site.Closed,
					})
				}
			}
		}

		body.Sitematrix = nil

		return queue, nil
	}
}
