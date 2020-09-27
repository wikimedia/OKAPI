package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"okapi/helpers/download"
	"okapi/helpers/exception"
	user_helper "okapi/helpers/user"
	"okapi/models/roles"
)

// DownloadRestrictions validate numbe of user downloads
func DownloadRestrictions(quantities map[roles.Type]int) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			user, err := user_helper.FromContext(c)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, exception.Message(fmt.Errorf("You need to be logged in")))
				return
			}

			quantity, exists := quantities[user.RoleID]

			if !exists {
				c.AbortWithStatusJSON(http.StatusForbidden, exception.Message(fmt.Errorf("Not enough permissions")))
				return
			}

			if quantity == -1 {
				c.Next()
				return
			}

			var count int

			counter := download.NewCounter(user)
			count, err = counter.Get()

			if err != nil && err != redis.Nil {
				c.AbortWithStatusJSON(http.StatusForbidden, exception.Message(err))
				return
			}

			if err == redis.Nil {
				counter.Set(0)

				count = 0
			}

			if quantity <= count {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					exception.Message(fmt.Errorf("You have reached a number of project downloads")),
				)
				return
			}
		}
	}
}
