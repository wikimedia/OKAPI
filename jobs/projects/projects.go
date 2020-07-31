package projects

import (
	"encoding/json"

	"okapi/helpers/logger"

	"okapi/lib/task"
	"okapi/lib/wiki"
)

// Task getting all projects from sitematrix
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	client := wiki.Client("https://en.wikipedia.org/")
	res, err := client.R().SetFormData(map[string]string{
		"action":        "sitematrix",
		"format":        "json",
		"formatversion": "2",
	}).Post("w/api.php")

	if err != nil {
		logger.JOB.Panic(logger.Message{
			ShortMessage: "Job: 'projects' exec failed",
			FullMessage:  err.Error(),
		})
	}

	body := wiki.Projects{}
	json.Unmarshal(res.Body(), &body)

	return Pool(&body), Worker, nil
}
