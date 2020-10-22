package pull

import "okapi/models"

func updatePage(page *models.Page, fields ...string) error {
	_, err := models.DB().Model(page).Column(fields...).Where("id = ? and project_id = ?", page.ID, page.ProjectID).Update()
	return err
}
