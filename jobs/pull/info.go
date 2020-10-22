package pull

import (
	"okapi/models"
)

func getInfo(page *models.Page) map[string]interface{} {
	return map[string]interface{}{
		"_title": page.Title,
		"_id":    page.ID,
	}
}
