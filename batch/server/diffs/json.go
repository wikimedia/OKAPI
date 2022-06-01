package diffs

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "okapi-diffs/server/diffs/protos"
)

// ErrExportStorage type issue with local storage
var ErrExportStorage = errors.New("local storage should work with os.File")

// JSON put all files content into one file and add to writer
func JSON(tarbal *tar.Writer, files chan exportFile, store *ExportStorage, req *pb.ExportRequest, res *pb.ExportResponse) error {
	path := fmt.Sprintf("%s/%s", store.Tmp, fmt.Sprintf("%s.ndjson", req.DbName))
	storef, err := store.Local.Create(path)

	if err != nil {
		return err
	}

	if storef == nil {
		return ErrExportFileIsNil
	}

	var tmpf *os.File

	if val, ok := storef.(*os.File); ok {
		tmpf = val
	} else {
		return ErrExportStorage
	}

	defer tmpf.Close() // #nosec G307

	for file := range files {
		res.Total++

		if _, err := tmpf.Write(file.Buffer.Bytes()); err != nil {
			res.Errors++
			log.Println(err)
			continue
		}

		if _, err := tmpf.WriteString("\n"); err != nil {
			res.Errors++
			log.Println(err)
		}
	}

	stat, err := store.Local.Stat(path)

	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    fmt.Sprintf("%s.ndjson", req.DbName),
		Size:    stat.Size(),
		Mode:    0766,
		ModTime: time.Now().UTC(),
	}

	if err := tarbal.WriteHeader(header); err != nil {
		return err
	}

	jsonf, err := store.Local.Get(path)

	if err != nil {
		return err
	}

	defer jsonf.Close()

	if _, err := io.Copy(tarbal, jsonf); err != nil {
		return err
	}

	return store.Local.Delete(path)
}
