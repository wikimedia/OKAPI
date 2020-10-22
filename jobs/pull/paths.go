package pull

import (
	"okapi/models"
	"strings"
)

func getHTMLPath(page *models.Page) string {
	return "html/" + strings.Replace(page.SiteURL, "https://", "", 1) + "/" + page.Title + ".html"
}

func getWikitextPath(page *models.Page) string {
	return "wikitext/" + strings.Replace(page.SiteURL, "https://", "", 1) + "/" + page.Title + ".wikitext"
}
