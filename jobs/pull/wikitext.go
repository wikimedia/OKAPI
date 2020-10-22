package pull

import (
	"bytes"
	"fmt"
	"net/http"
	"okapi/lib/storage"
	"okapi/lib/wiki"
	"okapi/models"
)

func getWikitext(page *models.Page) error {
	wikitext, status, err := wiki.Client(page.SiteURL).GetRevisionWikitext(page.Title, page.Revision)

	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("status code: %d", status)
	}

	return storage.Local.Client().Put(page.WikitextPath, bytes.NewReader(wikitext))
}
