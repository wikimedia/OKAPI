package dump

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Folder get latest dumps folder
func Folder() (string, error) {
	res, err := Client().R().Get("/other/pagetitles/")

	if err != nil {
		return "", err
	}

	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("status code: %d", res.StatusCode())
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		return "", err
	}

	folder := ""
	doc.Find("pre a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		if name != "../" && len(folder) <= 0 {
			folder = strings.TrimSuffix(name, "/")
		}
	})

	if len(folder) <= 0 {
		return "", fmt.Errorf("can't find the folder")
	}

	return folder, nil
}
