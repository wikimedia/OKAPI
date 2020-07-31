package routes

import (
	"net/http"
	"time"

	"okapi/helpers/exception"
	"okapi/lib/storage"

	"github.com/gin-gonic/gin"
)

// URL gets one-time use link to download dump from S3
func URL(c *gin.Context) {
	path := c.Query("path")
	store := storage.Remote.Client()
	url, err := store.Link(path, 1*time.Minute)
	if err != nil {
		c.JSON(http.StatusNotFound, exception.Message(err))
		return
	}

	c.String(http.StatusOK, url)
}
