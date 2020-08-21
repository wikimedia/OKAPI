package routes

import (
	"net/http"
	"okapi/lib/env"
	"okapi/lib/test_api"
	"testing"
)

func TestDelete(t *testing.T) {
	env.Context.Fill()
	res, err := testApi.Client().R().Delete("/example/1")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusNoContent {
		t.Error(res.Status())
	}
}
