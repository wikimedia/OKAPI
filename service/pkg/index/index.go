package index

import "time"

// Avaliable indexes
const (
	Project = "project"
	Page    = "page"
)

// DocProject index struct
type DocProject struct {
	ID            int       `json:"-"`
	DbName        string    `json:"db_name"`
	SiteName      string    `json:"site_name"`
	SiteCode      string    `json:"site_code"`
	SiteURL       string    `json:"site_url"`
	Lang          string    `json:"lang"`
	LangName      string    `json:"lang_name"`
	LangLocalName string    `json:"lang_local_name"`
	Active        bool      `json:"active"`
	Size          float64   `json:"size"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DocPage index struct
type DocPage struct {
	ID            int       `json:"-"`
	Title         string    `json:"title"`
	Name          string    `json:"name"`
	NsID          int       `json:"ns_id"`
	DbName        string    `json:"db_name"`
	Lang          string    `json:"lang"`
	LangName      string    `json:"lang_name"`
	LangLocalName string    `json:"lang_local_name"`
	SiteCode      string    `json:"site_code"`
	SiteURL       string    `json:"site_url"`
	UpdatedAt     time.Time `json:"updated_at"`
}
