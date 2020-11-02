package pull

import (
	"encoding/json"
	"fmt"

	"okapi/jobs/pull"
	"okapi/models"
)

// Worker scan queue processor
func Worker(payload string) (string, map[string]interface{}, error) {
	pages := []*models.Page{}
	err := json.Unmarshal([]byte(payload), &pages)
	workerInfo := map[string]interface{}{
		"_worker": "pull",
	}

	if err != nil {
		return "", workerInfo, err
	}

	if len(pages) <= 0 {
		return "", workerInfo, fmt.Errorf("index out of range")
	}

	message, info, err := pull.Worker(1, pages[0])

	for name, value := range workerInfo {
		info[name] = value
	}

	return message, info, err
}
