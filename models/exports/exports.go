package exports

import "okapi/models"

// Resource export resource type
type Resource string

// Export type names
const (
	Project Resource = "project"
	Page    Resource = "page"
)

// Types export types models map
var Types = map[Resource]interface{}{
	Project: &models.Project{},
	Page:    &models.Page{},
}
