package minifier

import (
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

var mini *minify.M

// Client get minify client
func Client() *minify.M {
	if mini != nil {
		return mini
	}

	mini := minify.New()
	mini.AddFunc("text/html", html.Minify)
	return mini
}
