package export

import "bytes"

type file struct {
	name   string
	buffer *bytes.Buffer
}
