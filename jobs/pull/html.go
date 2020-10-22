package pull

import (
	"bytes"
	"fmt"
	"net/http"
	"okapi/lib/storage"
	"okapi/lib/wiki"
	"okapi/models"
)

func getHTML(page *models.Page) error {
	html, status, err := wiki.Client(page.SiteURL).GetRevisionHTML(page.Title, page.Revision)

	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("status code: %d", status)
	}

	return storage.Local.Client().Put(page.Path, bytes.NewReader(html))
}
