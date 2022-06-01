package schema

// Protection level for the page
type Protection struct {
	Type   string `json:"type,omitempty"`
	Level  string `json:"level,omitempty"`
	Expiry string `json:"expiry,omitempty"`
}
