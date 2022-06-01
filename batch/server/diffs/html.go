package diffs

import (
	"archive/tar"
	"io"
	"log"
	"time"

	pb "okapi-diffs/server/diffs/protos"
)

// HTML add files to tar writer
func HTML(tarbal *tar.Writer, files chan exportFile, res *pb.ExportResponse) error {
	for file := range files {
		res.Total++

		err := tarbal.WriteHeader(&tar.Header{
			Name:    file.Title,
			Size:    int64(file.Buffer.Len()),
			Mode:    0766,
			ModTime: time.Now().UTC(),
		})

		if err != nil {
			res.Errors++
			log.Println(err)
			continue
		}

		if _, err := io.Copy(tarbal, file.Buffer); err != nil {
			res.Errors++
			log.Println(err)
		}
	}

	return nil
}
