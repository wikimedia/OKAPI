package scan

import (
	"encoding/json"
	"fmt"
	"okapi/lib/queue"

	"okapi/jobs/scan"
	"okapi/models"
)

// Worker scan queue processor
func Worker(payload string) (string, map[string]interface{}, error) {
	pages := []*models.Page{}
	err := json.Unmarshal([]byte(payload), &pages)
	workerInfo := map[string]interface{}{
		"_worker": "scan",
	}

	if err != nil {
		return "", workerInfo, err
	}

	if len(pages) <= 0 {
		return "", workerInfo, fmt.Errorf("Index out of range")
	}

	message, info, err := scan.Worker(1, pages[0])

	for name, value := range workerInfo {
		info[name] = value
	}

	if err == nil {
		queue.Sync.Add(pages[0])
	}

	return message, info, err
}
