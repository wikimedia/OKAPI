package schema

const (
	NamespaceArticle  = 0
	NamespaceFile     = 6
	NamespaceCategory = 14
	NamespaceTemplate = 10
)

// Namespace schema
type Namespace struct {
	Name       string `json:"name"`
	Identifier int    `json:"identifier"`
}
