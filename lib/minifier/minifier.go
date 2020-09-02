package minifier

import (
	"fmt"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

var mini *minify.M

// Client get minify client
func Client() *minify.M {
	if mini != nil {
		return mini
	}

	mini = minify.New()
	mini.AddFunc("text/html", html.Minify)
	return mini
}

// Init function to run on star
func Init() error {
	Client()

	if mini == nil {
		return fmt.Errorf("Minify client not available")
	}

	return nil
}
