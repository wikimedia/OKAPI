package revision

import (
	"encoding/json"
	"okapi/lib/ores"
	"okapi/lib/queue"
	"okapi/models"
)

// Payload for revision processor
type Payload struct {
	Page     models.Page
	Project  models.Project
	Revision int
}

// Worker revision queue processor
func Worker(rawPayload string) (string, map[string]interface{}, error) {
	payload := []*Payload{}
	err := json.Unmarshal([]byte(rawPayload), &payload)
	workerInfo := map[string]interface{}{
		"_worker":          "revision",
		"_page_title":      payload[0].Page.Title,
		"_project_db_name": payload[0].Project.DBName,
	}

	if err != nil {
		return "", workerInfo, err
	}

	page := payload[0].Page
	project := payload[0].Project
	page.SetRevision(payload[0].Revision)

	if ores.Damaging.CanScore(project.DBName) {
		return "", workerInfo, nil
	}

	if err != nil {
		return "", workerInfo, err
	}

	_, err = models.DB().
		Model(&page).
		Where("id = ? and project_id = ?", page.ID, project.ID).
		Update("revision", "updates", "scores", "updated_at")

	if err == nil {
		queue.PagePull.Add(page)
	}

	return "", workerInfo, err
}
