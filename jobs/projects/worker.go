package projects

import (
	"fmt"
	"strconv"

	"okapi/lib/task"
	"okapi/models"
)

// Worker projects unit processor
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	message := "name:'%s', id:#%d"
	project := payload.(*models.Project)
	info := map[string]interface{}{
		"_db_name":   project.DBName,
		"_lang_name": project.LangName,
	}

	models.DB().Model(project).Where("db_name = ?", project.DBName).Select()
	err := models.Save(project)

	if err != nil {
		err = fmt.Errorf(message+", %s", project.DBName, project.ID, err)
	}

	if err == nil {
		_, err = models.DB().Exec("CREATE TABLE IF NOT EXISTS  pages_" + project.DBName + " PARTITION OF pages FOR VALUES in ('" + strconv.Itoa(project.ID) + "')")
	}

	return fmt.Sprintf(message, project.DBName, project.ID), info, err
}
