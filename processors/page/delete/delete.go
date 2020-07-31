package deletePage

import (
	"encoding/json"
	"fmt"
	"okapi/lib/storage"
	"okapi/models"
)

// Worker delete-page processor
func Worker(payload string) (string, map[string]interface{}, error) {
	var pages []*models.Page
	message := "page title: '%s', page id: #%d"
	err := json.Unmarshal([]byte(payload), &pages)
	workerInfo := map[string]interface{}{
		"_worker": "delete-page",
	}

	if err != nil {
		return "", workerInfo, err
	}

	if len(pages) <= 0 {
		return "", workerInfo, fmt.Errorf("index out of range")
	}

	page := pages[0]

	err = storage.Local.Client().Delete(page.Path)
	if err != nil {
		return "", workerInfo, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	err = models.Delete(page)
	if err != nil {
		return "", workerInfo, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	return fmt.Sprintf(message, page.Title, page.ID), workerInfo, nil
}
