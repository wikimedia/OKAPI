package utils

import "fmt"

// SiteURL form site url from domain
func SiteURL(domain string) string {
	return fmt.Sprintf("https://%s", domain)
}
