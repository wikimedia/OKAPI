package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	user_helper "okapi/helpers/user"
	"okapi/models/permissions"
)

// PermissionsList list of permission types
type PermissionsList []permissions.Type

// Permissions middleware handler
func Permissions(types PermissionsList) func() gin.HandlerFunc {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			user, err := user_helper.FromContext(c)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You need to be logged in"})
			}

		OUTER:
			for _, permissionType := range types {
				for _, rolePermissions := range user.Role.Permissions {
					if permissionType == rolePermissions {
						continue OUTER
					}
				}

				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Not enough permissions"})

				return
			}
		}
	}
}
