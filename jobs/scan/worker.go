package scan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Worker processor for one title and store him into database as page
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	message := "page title: '%s', page id: #%d"
	page := payload.(*models.Page)
	info := map[string]interface{}{
		"_title":  page.Title,
		"_status": http.StatusOK,
		"_id":     page.ID,
	}

	res, err := wiki.Client(page.SiteURL).R().Get("/api/rest_v1/page/title/" + url.QueryEscape(page.Title))

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	if res.StatusCode() != http.StatusOK {
		info["_status"] = res.StatusCode()
		return "", info, fmt.Errorf(message+", status code: %d", page.Title, page.ID, res.StatusCode())
	}

	title := wiki.Title{}
	err = json.Unmarshal(res.Body(), &title)
	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	if title.Items[0].Redirect {
		info["_status"] = 301
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, "page is a redirect!")
	}

	wikiPage := title.Items[0]
	models.DB().Model(page).Column("id", "created_at").Where("title = ? and project_id = ?", page.Title, page.ProjectID).Select()
	page.TID = wikiPage.TID
	page.Revision = wikiPage.Revision
	page.Title = wikiPage.Title
	page.Lang = wikiPage.PageLanguage

	err = models.Save(page)
	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	info["_id"] = page.ID

	return fmt.Sprintf(message, page.Title, page.ID), info, nil
}
