package progress

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"gopkg.in/gookit/color.v1"
)

// Reader progress reader
type Reader struct {
	Message string
	Reader  io.Reader
	Num     int64
	Size    int64
	LogTime time.Time
}

func (reader *Reader) Read(p []byte) (int, error) {
	num, err := reader.Reader.Read(p)
	atomic.AddInt64(&reader.Num, int64(num))

	if !reader.LogTime.Add(1 * time.Second).Before(time.Now()) {
		return num, err
	}

	reader.LogTime = time.Now()

	if reader.Size != 0 {
		color.Blue.Println(fmt.Sprintf(reader.Message, int(float32(reader.Num*100)/float32(reader.Size))))
	}

	return num, err
}
