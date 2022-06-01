package v1

import (
	"okapi-streams/modules/v1/pages"

	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

// Init initialize v1 endpoints
func Init() []func() httpmod.Module {
	return []func() httpmod.Module{
		pages.Init,
	}
}
