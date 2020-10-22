package routes

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"okapi/helpers/exception"
	"okapi/lib/cache"
	"strconv"
)

// Get captcha ID
func Get(c *gin.Context) {
	response := struct {
		CaptchaID string `json:"captchaID"`
	}{captcha.New()}

	randomDigits := captcha.RandomDigits(6)

	var digits string

	for _, d := range randomDigits {
		digits += strconv.Itoa(int(d))
	}

	err := cache.Client().Set(response.CaptchaID, string(digits), captcha.Expiration).Err()

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
