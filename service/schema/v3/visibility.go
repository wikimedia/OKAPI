package schema

// Visibility representing visibility changes for parts of the revision
type Visibility struct {
	Text    bool `json:"text"`
	User    bool `json:"user"`
	Comment bool `json:"comment"`
}
