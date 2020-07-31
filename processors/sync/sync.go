package sync

import (
	"encoding/json"
	"fmt"

	"okapi/jobs/sync"
	"okapi/models"
)

// Worker scan queue processor
func Worker(payload string) (string, map[string]interface{}, error) {
	pages := []*models.Page{}
	err := json.Unmarshal([]byte(payload), &pages)
	workerInfo := map[string]interface{}{
		"_worker": "sync",
	}

	if err != nil {
		return "", workerInfo, err
	}

	if len(pages) <= 0 {
		return "", workerInfo, fmt.Errorf("Index out of range")
	}

	message, info, err := sync.Worker(1, pages[0])

	for name, value := range workerInfo {
		info[name] = value
	}

	return message, info, err
}
