package projects

import (
	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Pool function to get projects into the queue
func Pool(projects *wiki.Projects) func() ([]task.Payload, error) {
	return func() ([]task.Payload, error) {
		queue := []task.Payload{}

		if projects.Sitematrix != nil {
			for _, project := range projects.Sitematrix {
				for _, site := range project.Site {

					queue = append(queue, &models.Project{
						LangName:      project.Name,
						Lang:          project.Code,
						SiteName:      site.SiteName,
						SiteURL:       site.URL,
						SiteCode:      site.Code,
						DBName:        site.DBName,
						Dir:           project.Dir,
						LangLocalName: project.LocalName,
						Active:        !site.Closed,
					})
				}
			}
		}

		projects.Sitematrix = nil

		return queue, nil
	}
}
