package diffs

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/md5" // #nosec G501
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"okapi-diffs/schema/v3"
	pb "okapi-diffs/server/diffs/protos"
	"runtime"
	"sync"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/klauspost/pgzip"
)

// ErrExportFileIsNil singal that export returned nil (for stubs and new interfaces)
var ErrExportFileIsNil = errors.New("io.ReadWriteCloser is equal to nil")

// ExportStorage storage to get data from and to put archive into
type ExportStorage struct {
	Tmp      string // path to tmp directory
	Path     string // path to the files directory (where to get data from)
	MetaDest string // path to the meta destination directory (where to put the diff meta)
	Dest     string // path to the destination directory (where to put the dump)
	Local    interface {
		storage.Getter
		storage.Creator
		storage.Stater
		storage.Walker
		storage.Deleter
	}
	Remote storage.Putter
}

type exportFile struct {
	Title  string
	Buffer *bytes.Buffer
}

// Export generate new export file and upload it to the storage
func Export(ctx context.Context, req *pb.ExportRequest, store *ExportStorage, res *pb.ExportResponse) error {
	file, err := store.Local.Create(store.Dest)

	if err != nil {
		return err
	}

	if file == nil {
		return ErrExportFileIsNil
	}

	var proj *schema.Project

	write, read := new(sync.WaitGroup), new(sync.WaitGroup)
	files := make(chan exportFile, int(req.Workers))
	paths := make(chan string, int(req.Workers))
	gzip := pgzip.NewWriter(file)

	if err := gzip.SetConcurrency(1<<20, runtime.NumCPU()*2); err != nil {
		return err
	}

	tarbal := tar.NewWriter(gzip)

	// spin up workers that will be responsible for reading files from the filesystem
	read.Add(int(req.Workers))
	for i := 1; i <= int(req.Workers); i++ {
		go func() {
			defer read.Done()
			for path := range paths {
				file, err := store.Local.Get(path)

				if err != nil {
					log.Println(err)
					continue
				}

				data, err := ioutil.ReadAll(file)
				_ = file.Close()

				if err != nil {
					log.Println(err)
					continue
				}

				page := new(schema.Page)

				if err := json.Unmarshal(data, page); err != nil {
					log.Println(err)
					continue
				}

				if page.ArticleBody == nil || page.Namespace == nil {
					log.Printf("warning: empty schema payload: %v", page)
					continue
				}

				if proj == nil {
					proj = page.IsPartOf
					proj.InLanguage = page.InLanguage
					u, _ := url.Parse(page.URL)
					proj.URL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				}

				if page.Namespace.Identifier == int(req.Ns) {
					files <- exportFile{
						Title:  fmt.Sprintf("%s.json", page.Name),
						Buffer: bytes.NewBuffer(data),
					}
				}
			}
		}()
	}

	// create write worker
	write.Add(1)
	go func() {
		var err error

		defer write.Done()

		if err = JSON(tarbal, files, store, req, res); err != nil {
			log.Println(err)
		}
	}()

	err = store.Local.Walk(store.Path, func(path string) {
		paths <- path
	})

	close(paths)
	read.Wait()
	close(files)
	write.Wait()
	_ = tarbal.Close()
	_ = gzip.Close()
	_ = file.Close()

	if err != nil {
		return err
	}

	export, err := store.Local.Get(store.Dest)

	if err != nil {
		return err
	}

	defer func() {
		if err := store.Local.Delete(store.Dest); err != nil {
			log.Println(err)
		}
	}()
	defer export.Close()

	if err := store.Remote.Put(store.Dest, export); err != nil {
		return err
	}

	info, err := store.Local.Stat(store.Dest)

	if err != nil {
		return err
	}

	if proj == nil {
		return nil
	}

	if res.Total == 0 {
		return nil
	}

	// generate body md5 hash
	h := md5.New() // #nosec G401
	if _, err = io.Copy(h, export); err != nil {
		return err
	}

	size := math.Round((((float64)(info.Size())/1024)/1024)*100) / 100
	version := fmt.Sprintf("%x", h.Sum(nil))
	datetime := time.Now().UTC()
	meta := schema.Project{
		Name:         proj.Name,
		Identifier:   proj.Identifier,
		URL:          proj.URL,
		Version:      &version,
		DateModified: &datetime,
		Size: &schema.Size{
			Value:    size,
			UnitText: "MB",
		},
	}

	if proj.InLanguage != nil {
		meta.InLanguage = proj.InLanguage
	}

	metadata, err := json.Marshal(meta)

	if err != nil {
		return err
	}

	return store.Remote.Put(store.MetaDest, bytes.NewReader(metadata))
}
