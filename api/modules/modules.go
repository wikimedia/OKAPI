package modules

import (
	v1 "okapi-public-api/modules/v1"

	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

func Init() []func() httpmod.Module {
	return append([]func() httpmod.Module{}, v1.Init()...)
}
