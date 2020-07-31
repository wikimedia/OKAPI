package routes

import (
	"net/http"
	"okapi/lib/api"
	"okapi/lib/env"
	"testing"
)

func TestView(t *testing.T) {
	env.Context.Fill()
	res, err := api.Client().R().Get("/example/1")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusOK {
		t.Error(res.Status())
	}
}
