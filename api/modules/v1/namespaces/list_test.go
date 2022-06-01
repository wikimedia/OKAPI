package namespaces

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"okapi-public-api/schema/v3"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const listTestURL = "/v1/namespaces"

func createListTestServer() http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Handle(http.MethodGet, listTestURL, List())
	return router
}

func TestList(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(createListTestServer())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, listTestURL))
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	defer res.Body.Close()
	assert.NotEmpty(body)
	assert.True(json.Valid(body))

	namespaces := make([]schema.Namespace, 0)
	assert.NoError(json.Unmarshal(body, &namespaces))
	assert.IsType(schema.Namespace{}, namespaces[0])
}
