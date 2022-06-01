package contenttype

import "errors"

// ErrInvalidContentType when trying to query unavailable ContentType in elastic
var ErrInvalidContentType = errors.New("invalid search Content Type, available Content Types: 'html', 'json', 'wikitext'")

// ContentType available to query from elastic
type ContentType string

// Available ContentTypes for download
const (
	HTML  ContentType = "html"
	WText ContentType = "wikitext"
	JSON  ContentType = "json"
)

// Validate checks if Resource name is valid
func (r ContentType) Validate() error {
	switch r {
	case HTML, WText, JSON:
		return nil
	default:
		return ErrInvalidContentType
	}
}
