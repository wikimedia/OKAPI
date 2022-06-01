package utils

import "fmt"

// DateFormat for directories
const DateFormat = "2006-01-02"

// Format get formatted dir path
func Format(dir string, dbName string, contentType string, title string, fileType string) string {
	return fmt.Sprintf("page/%s/%s/%s/%s.%s", dir, dbName, contentType, title, fileType)
}
