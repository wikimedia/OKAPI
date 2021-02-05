package utils

import (
	"strings"
)

// Lang exract language from domain name
func Lang(domain string) string {
	parts := strings.Split(domain, ".")

	if len(parts) > 0 && len(parts[0]) > 0 {
		return parts[0]
	}

	return ""
}
