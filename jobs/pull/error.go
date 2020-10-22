package pull

import (
	"fmt"
	"okapi/models"
)

func getError(page *models.Page, err error) error {
	return fmt.Errorf(message+", %s", page.Title, page.ID, err)
}
