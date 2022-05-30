package schema

// Version page versions meta data
type Version struct {
	Identifier      int      `json:"identifier,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	IsMinorEdit     bool     `json:"is_minor_edit,omitempty"`
	IsFlaggedStable bool     `json:"is_flagged_stable,omitempty"`
	Scores          *Scores  `json:"scores,omitempty"`
	Editor          *Editor  `json:"editor,omitempty"`
}
