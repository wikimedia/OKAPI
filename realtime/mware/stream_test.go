package mware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const streamTestURL = "/stream"

func createStreamServer() http.Handler {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(Stream())
	router.Handle(http.MethodGet, streamTestURL, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return router
}

func TestStream(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(createStreamServer())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s%s", srv.URL, streamTestURL))
	assert.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	assert.Equal("text/event-stream; charset=utf-8", res.Header.Get("Content-Type"))
	assert.Equal("no-cache", res.Header.Get("Cache-Control"))
	assert.Equal("keep-alive", res.Header.Get("Connection"))
}
