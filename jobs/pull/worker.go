package pull

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"okapi/lib/storage"
	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Worker processing one page from the queue, getting html into s3
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	message := "page title: '%s', page id: #%d"
	page := payload.(*models.Page)
	page.Path = "html/" + strings.Replace(page.SiteURL, "https://", "", 1) + "/" + page.Title + ".html"
	info := map[string]interface{}{
		"_title":  page.Title,
		"_status": http.StatusOK,
		"_id":     page.ID,
	}

	html, status, err := wiki.Client(page.SiteURL).GetRevisionHTML(page.Title, page.Revision)

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	if status != http.StatusOK {
		info["_status"] = status
		return "", info, fmt.Errorf(message+", status code: %d", page.Title, page.ID, status)
	}

	if err != nil {
		return "", info, err
	}

	err = storage.Local.Client().Put(page.Path, bytes.NewReader(html))

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	_, err = models.DB().Model(page).Column("path").WherePK().Update()

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	return fmt.Sprintf(message, page.Title, page.ID), info, nil
}
