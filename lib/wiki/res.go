package wiki

import "time"

// Page API response for page
type Page struct {
	Title        string    `json:"title"`
	Revision     int       `json:"rev"`
	TID          string    `json:"tid"`
	Namespace    int       `json:"namespace"`
	UserID       int       `json:"user_id"`
	UserText     string    `json:"user_text"`
	Timestamp    time.Time `json:"timestamp"`
	Comment      string    `json:"comment"`
	PageLanguage string    `json:"page_language"`
	Redirect     bool      `json:"redirect"`
}

// Title API response or title
type Title struct {
	Items []*Page `json:"items"`
}

// Sitematrix wiki API payload response
type Sitematrix map[string]struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Dir       string `json:"dir"`
	LocalName string `json:"localname"`
	Site      []struct {
		URL      string `json:"url"`
		DBName   string `json:"dbname"`
		Code     string `json:"code"`
		SiteName string `json:"sitename"`
		Closed   bool   `json:"closed"`
	} `json:"site"`
}

// Projects wiki API req response body for projects
type Projects struct {
	Sitematrix Sitematrix `json:"sitematrix"`
}

// Revision wiki API response for revisions history
type Revision struct {
	RevID     int       `json:"revid"`
	ParentID  int       `json:"parentid"`
	Minor     bool      `json:"minor"`
	User      string    `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

// RevisionsHistory revisions history response
type RevisionsHistory struct {
	Query struct {
		Pages []struct {
			PageID    int        `json:"pageid"`
			Title     string     `json:"title"`
			Revisions []Revision `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

// RevisionWikitext wikitext by revision id
type RevisionWikitext struct {
	Query struct {
		Pages []struct {
			Revisions []struct {
				Slots struct {
					Main struct {
						Content string `json:"content"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

// Namespace single namespace
type Namespace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Namespaces namespaces list response
type Namespaces struct {
	Query struct {
		Namespaces map[int]Namespace
	}
}
