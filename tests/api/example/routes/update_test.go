package routes

import (
	"net/http"
	"okapi/lib/api"
	"okapi/lib/env"
	"testing"
)

func TestUpdate(t *testing.T) {
	env.Context.Fill()
	res, err := api.Client().R().Put("/example/1")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusOK {
		t.Error(res.Status())
	}
}
