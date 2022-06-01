package schema

import "time"

// Editor for the page revision
type Editor struct {
	Identifier  int        `json:"identifier,omitempty"`
	Name        string     `json:"name,omitempty"`
	EditCount   int        `json:"edit_count,omitempty"`
	Groups      []string   `json:"groups,omitempty"`
	IsBot       bool       `json:"is_bot,omitempty"`
	IsAnonymous bool       `json:"is_anonymous,omitempty"`
	DateStarted *time.Time `json:"date_started,omitempty"`
}
