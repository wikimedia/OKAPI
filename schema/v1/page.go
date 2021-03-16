package schema

import "time"

// PageKey key for page topics
type PageKey struct {
	Title  string `json:"title"`
	DbName string `json:"dbName"`
}

// PageURL represent structured content url object
type PageURL struct {
	Canonical string `json:"canonical"`
}

// PageBody struct to represent article body
type PageBody struct {
	HTML     string `json:"html"`
	Wikitext string `json:"wikitext"`
}

// Page content representation
type Page struct {
	Title        string    `json:"title"`
	PID          int       `json:"pid"`
	QID          string    `json:"qid,omitempty"`
	Revision     int       `json:"revision"`
	DbName       string    `json:"dbName"`
	InLanguage   string    `json:"inLanguage"`
	URL          PageURL   `json:"url"`
	Visible      *bool     `json:"visible,omitempty"`
	DateModified time.Time `json:"dateModified"`
	ArticleBody  PageBody  `json:"articleBody"`
	License      []string  `json:"license"`
}

// SetHTML set html body
func (p *Page) SetHTML(html string) {
	p.ArticleBody.HTML = html
}

// SetWikitext set wikitext body
func (p *Page) SetWikitext(wikitext string) {
	p.ArticleBody.Wikitext = wikitext
}
