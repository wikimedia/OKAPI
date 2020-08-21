package routes

import (
	"net/http"
	"okapi/lib/env"
	"okapi/lib/test_api"
	"testing"
)

func TestUpdate(t *testing.T) {
	env.Context.Fill()
	res, err := testApi.Client().R().Put("/example/1")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusOK {
		t.Error(res.Status())
	}
}
