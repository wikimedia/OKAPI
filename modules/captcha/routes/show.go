package routes

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"okapi/helpers/exception"
	"okapi/lib/cache"
	"strconv"
)

// Captcha properties
const (
	stdWidth  int = 250
	stdHeight int = 100
)

// Show return captcha image
func Show(c *gin.Context) {
	captchaID := c.Param("id")
	captchaSolution, err := cache.Client().Get(captchaID).Result()

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(fmt.Errorf("Captcha id is not valid")))
		return
	}

	digits := make([]byte, len(captchaSolution))

	for index, item := range captchaSolution {
		v, _ := strconv.Atoi(string(item))

		digits[index] = byte(v)
	}

	var content bytes.Buffer

	_, err = captcha.NewImage(captchaID, digits, stdWidth, stdHeight).WriteTo(&content)

	if err != nil {
		c.JSON(http.StatusBadRequest, exception.Message(err))
		return
	}

	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Write(content.Bytes())
}
