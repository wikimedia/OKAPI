package schema

import "time"

// Project schema
type Project struct {
	Name         string     `json:"name"`
	Identifier   string     `json:"identifier"`
	URL          string     `json:"url,omitempty"`
	Version      *string    `json:"version,omitempty"`
	DateModified *time.Time `json:"dateModified,omitempty"`
	InLanguage   *Language  `json:"inLanguage,omitempty"`
	Size         *Size      `json:"size,omitempty"`
}
