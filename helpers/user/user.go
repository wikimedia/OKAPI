package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"okapi/models"
)

// FromContext return a user from context
func FromContext(c *gin.Context) (*models.User, error) {
	userRaw, exists := c.Get("user")

	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	switch user := userRaw.(type) {
	case *models.User:
		return user, nil
	default:
		return nil, fmt.Errorf("user type is unknown")
	}
}
