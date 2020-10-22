package score

import (
	"encoding/json"
	"fmt"
	"okapi/lib/ores"
	"okapi/lib/queue"
	"okapi/models"
)

// Payload for revision processor
type Payload struct {
	Page     models.Page
	Project  models.Project
	Revision int
	Scores   map[string]ores.Stream
}

// Worker revision queue processor
func Worker(rawPayload string) (string, map[string]interface{}, error) {
	payload := []*Payload{}
	err := json.Unmarshal([]byte(rawPayload), &payload)
	workerInfo := map[string]interface{}{
		"_worker":          "score",
		"_page_title":      payload[0].Page.Title,
		"_project_db_name": payload[0].Project.DBName,
	}

	if err != nil {
		return "", workerInfo, err
	}

	page := payload[0].Page
	project := payload[0].Project
	page.SetRevision(payload[0].Revision)

	if stream, ok := payload[0].Scores["damaging"]; ok {
		err = scoreRevision(&page, &project, &stream.Probability)
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

func scoreRevision(page *models.Page, project *models.Project, probability *ores.Probability) error {
	threshold := project.GetThreshold(ores.Damaging)

	if threshold == nil {
		return fmt.Errorf("threshold model does not exist: project_id: %d, model: %s", project.ID, ores.Damaging)
	}

	if probability.False >= *threshold {
		page.SetScore(page.Revision, ores.Damaging, ores.Score{
			Prediction:  false,
			Probability: *probability,
		})
		return nil
	}

	return fmt.Errorf("the page revision is damaged: title %s; rev_id: %d", page.Title, page.Revision)
}
