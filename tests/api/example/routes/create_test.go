package routes

import (
	"net/http"
	"okapi/lib/api"
	"okapi/lib/env"
	"testing"
)

func TestCreate(t *testing.T) {
	env.Context.Fill()
	res, err := api.Client().R().Post("/example")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusCreated {
		t.Error(res.Status())
	}
}
