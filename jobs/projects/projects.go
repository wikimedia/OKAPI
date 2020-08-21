package projects

import (
	"okapi/helpers/logger"

	"okapi/lib/task"
	"okapi/lib/wiki"
)

// Task getting all projects from sitematrix
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	projects, _, err := wiki.Client("https://en.wikipedia.org/").GetSitematrix()

	if err != nil {
		logger.JOB.Panic(logger.Message{
			ShortMessage: "Job: 'projects' exec failed",
			FullMessage:  err.Error(),
		})
	}

	return Pool(projects), Worker, nil
}
