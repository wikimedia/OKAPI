package export

import "bytes"

type file struct {
	name   string
	path   string
	buffer *bytes.Buffer
}
