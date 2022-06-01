package schema

import "time"

// PageKey key for page topics inside kafka
type PageKey struct {
	Name     string `json:"name"`
	IsPartOf string `json:"isPartOf"`
}

// ArticleBody content of the page
type ArticleBody struct {
	HTML     string `json:"html"`
	Wikitext string `json:"wikitext"`
}

// PageVisibility representing visibility changes for parts of the page
type PageVisibility struct {
	Text    bool `json:"text"`
	User    bool `json:"user"`
	Comment bool `json:"comment"`
}

// Page schema
type Page struct {
	Name         string          `json:"name"`
	Identifier   int             `json:"identifier"`
	Version      int             `json:"version"`
	DateModified time.Time       `json:"dateModified,omitempty"`
	URL          string          `json:"url"`
	Namespace    *Namespace      `json:"namespace,omitempty"`
	InLanguage   *Language       `json:"inLanguage,omitempty"`
	MainEntity   *MainEntity     `json:"mainEntity,omitempty"`
	IsPartOf     *Project        `json:"isPartOf,omitempty"`
	ArticleBody  *ArticleBody    `json:"articleBody,omitempty"`
	License      []*License      `json:"license,omitempty"`
	Visibility   *PageVisibility `json:"visibility,omitempty"`
}

// SetHTML set html body
func (p *Page) SetHTML(html string) {
	p.ArticleBody.HTML = html
}

// SetWikitext set wikitext body
func (p *Page) SetWikitext(wikitext string) {
	p.ArticleBody.Wikitext = wikitext
}
