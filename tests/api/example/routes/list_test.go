package routes

import (
	"net/http"
	"okapi/lib/env"
	"okapi/lib/test_api"
	"testing"
)

func TestList(t *testing.T) {
	env.Context.Fill()
	res, err := testApi.Client().R().Get("/example")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusOK {
		t.Error(res.Status())
	}
}
