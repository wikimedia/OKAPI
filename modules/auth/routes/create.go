package routes

import (
	"fmt"
	"net/http"
	"okapi/lib/cache"

	"okapi/helpers/exception"
	"okapi/helpers/password"
	"okapi/models"

	"github.com/gin-gonic/gin"
)

type userParams struct {
	Username             string `form:"username" binding:"required,max=255"`
	Email                string `form:"email" binding:"required,email"`
	Password             string `form:"password" binding:"required,min=8"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required_with=Password,eqfield=Password"`
	CaptchaID            string `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaSolution      string `form:"captcha_solution" json:"captcha_solution" binding:"required"`
}

// Create a user handler
func Create(c *gin.Context) {
	var params userParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	trueSolution, _ := cache.Client().Get(params.CaptchaID).Result()

	cache.Client().Del(params.CaptchaID)

	if trueSolution != params.CaptchaSolution {
		c.JSON(http.StatusBadRequest, exception.Message(fmt.Errorf("Captcha is not valid")))
		return
	}

	hash, err := password.Hash(params.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	user := models.User{
		Email:    params.Email,
		Username: params.Username,
		Password: hash,
	}

	if err := models.Save(&user); err != nil {
		c.JSON(http.StatusForbidden, exception.Message(err))
	} else {
		c.JSON(http.StatusOK, &user)
	}
}
