package schema

// Project structured content schema
type Project struct {
	Name       string `json:"name"`
	DbName     string `json:"dbName"`
	InLanguage string `json:"inLanguage"`
	Size       string `json:"size"`
	URL        string `json:"URL"`
}
