package content

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"okapi-data-service/models"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/protsack-stephan/mediawiki-api-client"
)

// License default license for the pages
const License = "CC BY-SA"

// Storage storage options for different content types
type Storage struct {
	JSON  storage.Putter
	HTML  storage.Putter
	WText storage.Putter
}

// Structured content schema for json
type Structured struct {
	Title      string    `json:"title"`
	DbName     string    `json:"db_name"`
	PID        int       `json:"pid"`
	QID        string    `json:"qid"`
	URL        string    `json:"url"`
	Lang       string    `json:"lang"`
	Revision   int       `json:"revision"`
	RevisionDt time.Time `json:"revision_dt"`
	License    []string  `json:"license"`
	HTML       string    `json:"html"`
	Wikitext   string    `json:"wikitext"`
}

// Pull get and parse all the page data
func Pull(ctx context.Context, page *models.Page, storage *Storage, mwiki *mediawiki.Client) error {
	html := make(chan []byte, 1)
	wt := make(chan []byte, 1)
	errs := make(chan error, 2)

	go func() {
		data, err := mwiki.PageHTML(ctx, page.Title, page.Revision)

		if err != nil {
			errs <- err
			return
		}

		page.HTMLPath = fmt.Sprintf("html/%s/%s.html", page.DbName, page.Title)
		errs <- storage.HTML.Put(page.HTMLPath, bytes.NewReader(data))
		html <- data
	}()

	go func() {
		data, err := mwiki.PageWikitext(ctx, page.Title, page.Revision)

		if err != nil {
			errs <- err
			return
		}

		page.WikitextPath = fmt.Sprintf("wikitext/%s/%s.wt", page.DbName, page.Title)
		errs <- storage.WText.Put(page.WikitextPath, bytes.NewReader(data))
		wt <- data
	}()

	for i := 0; i < 2; i++ {
		err := <-errs

		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(Structured{
		Title:      page.Title,
		DbName:     page.DbName,
		PID:        page.PID,
		QID:        page.QID,
		URL:        fmt.Sprintf("%s/wiki/%s", page.SiteURL, page.Title),
		Lang:       page.Lang,
		Revision:   page.Revision,
		RevisionDt: page.RevisionDt,
		License:    []string{License},
		HTML:       string(<-html),
		Wikitext:   string(<-wt),
	})

	if err != nil {
		return err
	}

	page.JSONPath = fmt.Sprintf("json/%s/%s.json", page.DbName, page.Title)

	return storage.JSON.Put(page.JSONPath, bytes.NewReader(data))
}
