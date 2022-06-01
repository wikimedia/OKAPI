package v1

import (
	"okapi-public-api/modules/v1/diffs"
	"okapi-public-api/modules/v1/exports"
	"okapi-public-api/modules/v1/namespaces"
	"okapi-public-api/modules/v1/pages"
	"okapi-public-api/modules/v1/projects"

	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

// Init initialize v1 endpoints
func Init() []func() httpmod.Module {
	return []func() httpmod.Module{
		projects.Init,
		exports.Init,
		pages.Init,
		diffs.Init,
		namespaces.Init,
	}
}
