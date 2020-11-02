package jobs

import (
	"fmt"
	"okapi/helpers/state"
	"okapi/jobs"
	"okapi/lib/task"
	"okapi/models"
	"time"
)

func getState(name string, dbName string) *state.Client {
	return state.New(name+"_"+dbName, 24*time.Hour)
}

func getTask(name string) (job task.Task, err error) {
	job, exists := jobs.Tasks[task.Name(name)]

	if !exists {
		err = fmt.Errorf("task '%s' not found", string(name))
		return
	}

	return
}

func getProject(dbName string) (project models.Project, err error) {
	if dbName != "*" && len(dbName) > 0 {
		err = models.DB().Model(&project).Where("db_name = ?", dbName).Select()
	}

	return
}
