package namespaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

// Init namespaces endpoint.
func Init() httpmod.Module {
	return httpmod.Module{
		Path:       "/v1/namespaces",
		Middleware: []gin.HandlerFunc{},
		Routes: []httpmod.Route{
			{
				Method:  http.MethodGet,
				Handler: List(),
			},
		},
	}
}
