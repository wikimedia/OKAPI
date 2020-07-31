package routes

import (
	"net/http"

	"okapi/helpers/exception"
	"okapi/helpers/password"
	"okapi/models"

	"github.com/gin-gonic/gin"
)

type userParams struct {
	Email                string `form:"email" binding:"required,email"`
	Password             string `form:"password" binding:"required,min=8"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required_with=Password,eqfield=Password"`
}

// Create a user handler
func Create(c *gin.Context) {
	var params userParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	hash, err := password.Hash(params.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	user := models.User{
		Email:    params.Email,
		Password: hash,
	}

	if err := models.Save(&user); err != nil {
		c.JSON(http.StatusForbidden, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, &user)
	}
}
