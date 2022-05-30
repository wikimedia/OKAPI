package schema

import "time"

// Project schema
type Project struct {
	Name         string     `json:"name"`
	Identifier   string     `json:"identifier"`
	URL          string     `json:"url,omitempty"`
	Version      *string    `json:"version,omitempty"`
	DateModified *time.Time `json:"date_modified,omitempty"`
	InLanguage   *Language  `json:"in_language,omitempty"`
	Size         *Size      `json:"size,omitempty"`
}
