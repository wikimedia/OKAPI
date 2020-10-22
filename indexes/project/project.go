package project

import (
	"encoding/json"
	"okapi/lib/elastic"
	"strconv"
	"time"
)

// Name document collection name
const Name string = "project"

// Index search model for project
type Index struct {
	ID            int       `json:"-"`
	DBName        string    `json:"db_name"`
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
