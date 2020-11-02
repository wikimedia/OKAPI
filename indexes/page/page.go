package page

import (
	"encoding/json"
	"okapi/lib/elastic"
	"strconv"
	"time"
)

// Name document collection name
const Name string = "page"

// Index search model for page
type Index struct {
	ID            int       `json:"-"`
	NsID          int       `json:"ns_id"`
	Title         string    `json:"title"`
	SiteCode      string    `json:"site_code"`
	SiteURL       string    `json:"site_url"`
	Lang          string    `json:"lang"`
	LangName      string    `json:"lang_name"`
	LangLocalName string    `json:"lang_local_name"`
	ProjectID     int       `json:"project_id"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Update index inside elastic search
func (index *Index) Update() error {
	body, err := json.Marshal(index)

	if err == nil {
		elastic.Sync(strconv.Itoa(index.ID), Name, body)
	}

	return err
}

// Delete remove the index from elastic
func Delete(id int) {
	elastic.Delete(strconv.Itoa(id), Name)
}
