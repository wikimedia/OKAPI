package schema

import "time"

// PageKey key for page topics inside kafka
type PageKey struct {
	Name     string `json:"name"`
	IsPartOf string `json:"is_part_of"`
}

// ArticleBody content of the page
type ArticleBody struct {
	HTML     string `json:"html"`
	Wikitext string `json:"wikitext"`
}

// Page schema
type Page struct {
	Name               string        `json:"name"`
	Identifier         int           `json:"identifier,omitempty"`
	DateModified       *time.Time    `json:"date_modified,omitempty"`
	Protection         []*Protection `json:"protection,omitempty"`
	Version            *Version      `json:"version,omitempty"`
	URL                string        `json:"url,omitempty"`
	Namespace          *Namespace    `json:"namespace,omitempty"`
	InLanguage         *Language     `json:"in_language,omitempty"`
	MainEntity         *Entity       `json:"main_entity,omitempty"`
	AdditionalEntities []*Entity     `json:"additional_entities,omitempty"`
	Categories         []*Page       `json:"categories,omitempty"`
	Templates          []*Page       `json:"templates,omitempty"`
	Redirects          []*Page       `json:"redirects,omitempty"`
	IsPartOf           *Project      `json:"is_part_of,omitempty"`
	ArticleBody        *ArticleBody  `json:"article_body,omitempty"`
	License            []*License    `json:"license,omitempty"`
	Visibility         *Visibility   `json:"visibility,omitempty"`
}

// SetHTML set html body
func (p *Page) SetHTML(html string) {
	p.ArticleBody.HTML = html
}

// SetWikitext set wikitext body
func (p *Page) SetWikitext(wikitext string) {
	p.ArticleBody.Wikitext = wikitext
}
