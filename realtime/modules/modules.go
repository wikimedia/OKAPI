package modules

import (
	v1 "okapi-streams/modules/v1"

	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

// Init initialize all the modules
func Init() []func() httpmod.Module {
	return append([]func() httpmod.Module{}, v1.Init()...)
}
