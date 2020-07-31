package routes

import (
	"net/http"
	"okapi/lib/api"
	"okapi/lib/env"
	"testing"
)

func TestDelete(t *testing.T) {
	env.Context.Fill()
	res, err := api.Client().R().Delete("/example/1")

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode() != http.StatusNoContent {
		t.Error(res.Status())
	}
}
