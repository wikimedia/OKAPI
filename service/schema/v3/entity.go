package schema

// Entity schema for wikidata item
type Entity struct {
	Identifier string   `json:"identifier,omitempty"`
	URL        string   `json:"url,omitempty"`
	Aspects    []string `json:"aspects,omitempty"`
}
