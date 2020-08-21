package routes

import (
	"net/http"
	"okapi/lib/env"
	"okapi/lib/test_api"
	"testing"
)

func TestCreate(t *testing.T) {
	env.Context.Fill()
	res, err := testApi.Client().R().Post("/example")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusCreated {
		t.Error(res.Status())
	}
}
