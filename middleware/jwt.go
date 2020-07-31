package middleware

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"okapi/helpers/password"
	"okapi/lib/env"
	"okapi/models"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	ErrMissingLoginValues   = errors.New("missing Email or Password")
	ErrFailedAuthentication = errors.New("incorrect Email or Password")
)

const (
	Week        = (time.Hour * 24) * 7
	IdentityKey = "id"
)

var middleware *jwt.GinJWTMiddleware

func client() error {
	var err error

	middleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "UTC +0",
		Key:             []byte(env.Context.AuthSecretKey),
		Timeout:         Week,
		IdentityKey:     IdentityKey,
		PayloadFunc:     payloadFunc,
		IdentityHandler: identityHandler,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		Unauthorized:    unauthorized,
	})

	return err
}

// JWT middleware
func JWT() *jwt.GinJWTMiddleware {
	if middleware == nil {
		client()
	}

	return middleware
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(models.User); ok {
		return jwt.MapClaims{
			IdentityKey: user.ID,
		}
	}

	return jwt.MapClaims{
		IdentityKey: 0,
	}
}

func identityHandler(c *gin.Context) interface{} {
	var user models.User
	claims := jwt.ExtractClaims(c)

	models.DB().Model(&user).Where("id = ?", claims[IdentityKey]).Select()

	return &user
}

func authenticator(c *gin.Context) (interface{}, error) {
	var body login
	var user models.User

	if err := c.ShouldBind(&body); err != nil {
		return "", ErrMissingLoginValues
	}

	models.DB().Model(&user).Where("email = ?", body.Email).Select()

	if user.ID == 0 {
		return "", ErrFailedAuthentication
	}

	if err := password.CompareHashWithPassword(user.Password, body.Password); err != nil {
		return "", ErrFailedAuthentication
	}

	return user, nil
}

func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*models.User); ok && v.ID > 0 {
		return true
	}

	return false
}

func unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
