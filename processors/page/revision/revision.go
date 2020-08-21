package revision

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
}

// Worker revision queue processor
func Worker(rawPayload string) (string, map[string]interface{}, error) {
	payload := []*Payload{}
	err := json.Unmarshal([]byte(rawPayload), &payload)
	workerInfo := map[string]interface{}{
		"_worker": "revision",
	}

	if err != nil {
		return "", workerInfo, err
	}

	page := payload[0].Page
	project := payload[0].Project
	page.SetRevision(payload[0].Revision)
	score, scoreErr := ores.Damaging.ScoreOne(project.DBName, page.Revision)

	if scoreErr == nil {
		err = scoreRevision(&page, &project, score)
	}

	if err != nil {
		return "", workerInfo, err
	}

	err = models.Save(&page)

	if err == nil {
		queue.PagePull.Add(page)
	}

	return "", workerInfo, err
}

func scoreRevision(page *models.Page, project *models.Project, score *ores.Score) error {
	threshold := project.GetThreshold(ores.Damaging)

	if threshold == nil {
		return fmt.Errorf("threshold model does not exist: project_id: %d, model: %s", project.ID, ores.Damaging)
	}

	if score.Probability.False >= *threshold {
		page.SetScore(page.Revision, ores.Damaging, *score)
		return nil
	}

	return fmt.Errorf("the page revision is damaged: title %s; rev_id: %d", page.Title, page.Revision)
}
