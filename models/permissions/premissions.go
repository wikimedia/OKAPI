package permissions

// Type permissions type identifier
type Type string

// All available permissions
const (
	ProjectCreate Type = "project_create"
	ProjectDelete Type = "project_delete"
	ProjectBundle Type = "project_bundle"
)
