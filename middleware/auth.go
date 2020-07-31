package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"okapi/lib/env"
)

var users = gin.Accounts{}

// Auth middleware
func Auth() gin.HandlerFunc {
	if len(users) <= 0 {
		getAccounts()
	}

	return gin.BasicAuth(users)
}

// Accounts setting accounts string
func getAccounts() {
	accounts := strings.Split(env.Context.APIAuth, ",")

	for _, account := range accounts {
		cred := strings.Split(account, ":")
		users[cred[0]] = cred[1]
	}
}
