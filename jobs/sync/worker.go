package sync

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"okapi/lib/minifier"
	"okapi/lib/storage"
	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Worker processing one page from the queue, getting html into s3
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	message := "page title: '%s', page id: #%d"
	page := payload.(*models.Page)
	page.Path = "html/" + strings.Replace(page.SiteURL, "https:", "", 1) + "/" + page.Lang + "_" + page.Title + ".html"
	info := map[string]interface{}{
		"_title":  page.Title,
		"_status": http.StatusOK,
		"_id":     page.ID,
	}

	res, err := wiki.Client(page.SiteURL).R().Get("/api/rest_v1/page/html/" + url.QueryEscape(page.Title))

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	if res.StatusCode() != http.StatusOK {
		info["_status"] = res.StatusCode()
		return "", info, fmt.Errorf(message+", status code: %d", page.Title, page.ID, res.StatusCode())
	}

	html := string(res.Body())
	minified, err := minifier.Client().String("text/html", html)
	if err == nil {
		html = minified
	}

	err = storage.Local.Client().Put(page.Path, strings.NewReader(minified))
	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	err = models.Save(page)
	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	return fmt.Sprintf(message, page.Title, page.ID), info, nil
}
