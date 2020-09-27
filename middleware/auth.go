package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"okapi/lib/env"
)

var authUsers = gin.Accounts{}

// Auth middleware
func Auth() gin.HandlerFunc {
	if len(authUsers) <= 0 {
		getAccounts()
	}

	return gin.BasicAuth(authUsers)
}

// Accounts setting accounts string
func getAccounts() {
	accounts := strings.Split(env.Context.APIAuth, ",")

	for _, account := range accounts {
		cred := strings.Split(account, ":")
		authUsers[cred[0]] = cred[1]
	}
}
