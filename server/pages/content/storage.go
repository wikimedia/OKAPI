package content

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"okapi-data-service/models"
	"okapi-data-service/schema/v1"
	"strings"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
	"github.com/protsack-stephan/mediawiki-api-client"
)

type store interface {
	storage.Putter
	storage.Deleter
}

// Storage all available content sotrages
type Storage struct {
	HTML   store
	WText  store
	JSON   store
	Remote store
}

// Pull get page data and store inside the storage
func (s *Storage) Pull(ctx context.Context, page *models.Page, mwiki *mediawiki.Client) (*schema.Page, error) {
	html := make(chan string, 1)
	wt := make(chan string, 1)
	errs := make(chan error, 2)

	go func() {
		data, err := mwiki.PageHTML(ctx, page.Title, page.Revision)

		if err != nil {
			errs <- err
			return
		}

		mini := minify(string(data))
		page.HTMLPath = fmt.Sprintf("html/%s/%s.html", page.DbName, page.Title)
		errs <- s.HTML.Put(page.HTMLPath, strings.NewReader(mini))
		html <- mini
	}()

	go func() {
		data, err := mwiki.PageWikitext(ctx, page.Title, page.Revision)

		if err != nil {
			errs <- err
			return
		}

		mini := minify(string(data))
		page.WikitextPath = fmt.Sprintf("wikitext/%s/%s.wikitext", page.DbName, page.Title)
		errs <- s.WText.Put(page.WikitextPath, strings.NewReader(mini))
		wt <- mini
	}()

	for i := 0; i < 2; i++ {
		err := <-errs

		if err != nil {
			return nil, err
		}
	}

	cont := NewStructured(page)
	cont.SetHTML(<-html)
	cont.SetWikitext(<-wt)

	data, err := json.Marshal(cont)

	if err != nil {
		return nil, err
	}

	page.JSONPath = fmt.Sprintf("json/%s/%s.json", page.DbName, page.Title)

	go func() {
		errs <- s.JSON.Put(page.JSONPath, bytes.NewReader(data))
	}()

	go func() {
		errs <- s.Remote.Put(fmt.Sprintf("page/%s", page.JSONPath), bytes.NewReader(data))
	}()

	for i := 0; i < 2; i++ {
		err := <-errs

		if err != nil {
			return nil, err
		}
	}

	return cont, nil
}

// Delete remove page from all storages
func (s *Storage) Delete(_ context.Context, page *models.Page) error {
	var result error

	if err := s.HTML.Delete(page.HTMLPath); err != nil {
		result = err
	}

	if err := s.WText.Delete(page.WikitextPath); err != nil {
		result = err
	}

	if err := s.JSON.Delete(page.JSONPath); err != nil {
		result = err
	}

	if err := s.Remote.Delete(fmt.Sprintf("page/%s", page.JSONPath)); err != nil {
		result = err
	}

	return result
}
