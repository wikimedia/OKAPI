package content

import "strings"

func minify(data string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return -1
		}

		return r
	}, data)
}
