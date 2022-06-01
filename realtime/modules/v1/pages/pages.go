package pages

import (
	"net/http"
	"time"

	"okapi-streams/lib/env"
	"okapi-streams/mware"
	"okapi-streams/pkg/consumer"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/gin-toolkit/httpmod"
)

const topicPageUpdate = "aws.data-service.page-update.3"
const topicPageDelete = "aws.data-service.page-delete.3"
const topicPageVisibility = "aws.data-service.page-visibility.3"

// Update page update stream
// @Summary Get page update events
// @Tags pages
// @Description Returns stream of page structured data
// @ID v1-page-update
// @Security ApiKeyAuth
// @Param offset query number false "Offset"
// @Param since query string false "Since Date (in RFC3339 '2006-01-02T15:04:05Z07:00' or as a timestamp in milliseconds)"
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Failure 500 {object} httperr.Error
// @Router /v1/page-update [get]
func Update() gin.HandlerFunc {
	return Stream(topicPageUpdate, env.KafkaBroker, time.Second*1, consumer.NewConsumer)
}

// Delete page update stream
// @Summary Get page delete events
// @Tags pages
// @Description Returns stream of page delete events
// @ID v1-page-delete
// @Security ApiKeyAuth
// @Param offset query number false "Offset"
// @Param since query string false "Since Date (in RFC3339 '2006-01-02T15:04:05Z07:00' or as a timestamp in milliseconds)"
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Failure 500 {object} httperr.Error
// @Router /v1/page-delete [get]
func Delete() gin.HandlerFunc {
	return Stream(topicPageDelete, env.KafkaBroker, time.Minute*5, consumer.NewConsumer)
}

// Visibility page visibility change stream
// @Summary Get page visibility change events
// @Tags pages
// @Description Returns stream of page visibility change events
// @ID v1-page-visibility
// @Security ApiKeyAuth
// @Param offset query number false "Offset"
// @Param since query string false "Since Date (in RFC3339 '2006-01-02T15:04:05Z07:00' or as a timestamp in milliseconds)"
// @Failure 400 {object} httperr.Error
// @Failure 404 {object} httperr.Error
// @Failure 500 {object} httperr.Error
// @Router /v1/page-visibility [get]
func Visibility() gin.HandlerFunc {
	return Stream(topicPageVisibility, env.KafkaBroker, time.Hour*1, consumer.NewConsumer)
}

// Init for page endpoints
func Init() httpmod.Module {
	return httpmod.Module{
		Path: "/v1",
		Middleware: []gin.HandlerFunc{
			mware.Stream(),
		},
		Routes: []httpmod.Route{
			{
				Path:    "/page-update",
				Method:  http.MethodGet,
				Handler: Update(),
			},
			{
				Path:    "/page-delete",
				Method:  http.MethodGet,
				Handler: Delete(),
			},
			{
				Path:    "/page-visibility",
				Method:  http.MethodGet,
				Handler: Visibility(),
			},
		},
	}
}
