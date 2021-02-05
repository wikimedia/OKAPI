package pages

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"okapi-data-service/models"
	pb "okapi-data-service/server/pages/protos"
	"runtime"
	"sync"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/go-pg/pg/v10/orm"
	"github.com/klauspost/pgzip"
)

// ErrExportFileIsNil singal that export returned nil (for stubs and new interfaces)
var ErrExportFileIsNil = errors.New("io.ReadWriteCloser is equal to nil")

// ExportStorage storage to get data from and to put archive into
type ExportStorage struct {
	Loc  string // path to the files directory (where to get data from)
	Dest string // path to the destination directory (where to put the dump)
	To   interface {
		storage.Getter
		storage.Creator
		storage.Stater
	}
	From interface {
		storage.Getter
		storage.Walker
	}
	Remote storage.Putter
}

type exportRepo interface {
	repository.Finder
	repository.Updater
}

type exportFile struct {
	Title  string
	Buffer *bytes.Buffer
}

// Export generate new export file and upload it to the storage
func Export(ctx context.Context, req *pb.ExportReqest, repo exportRepo, store *ExportStorage) (*pb.ExportResponse, error) {
	proj := new(models.Project)
	err := repo.Find(ctx, proj, func(q *orm.Query) *orm.Query {
		return q.
			Where("db_name = ?", req.DbName)
	})

	if err != nil {
		return nil, err
	}

	file, err := store.To.Create(store.Dest)

	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, ErrExportFileIsNil
	}

	res := new(pb.ExportResponse)
	write, read := new(sync.WaitGroup), new(sync.WaitGroup)
	files := make(chan exportFile, int(req.Workers))
	paths := make(chan string, int(req.Workers))
	gzip := pgzip.NewWriter(file)

	if err := gzip.SetConcurrency(1<<20, runtime.NumCPU()*2); err != nil {
		return nil, err
	}

	tarbal := tar.NewWriter(gzip)

	// spin up workers that will be responsible for reading files from the filesystem
	read.Add(int(req.Workers))
	for i := 1; i <= int(req.Workers); i++ {
		go func() {
			defer read.Done()
			for path := range paths {
				file, err := store.From.Get(path)

				if err != nil {
					log.Println(err)
					return
				}

				data, err := ioutil.ReadAll(file)
				file.Close()

				if err != nil {
					log.Println(err)
					return
				}

				files <- exportFile{
					Title:  path[len(store.Loc)+1:],
					Buffer: bytes.NewBuffer(data),
				}
			}
		}()
	}

	// create write worker
	write.Add(1)
	go func() {
		defer write.Done()
		for file := range files {
			res.Total++

			err := tarbal.WriteHeader(&tar.Header{
				Name: file.Title,
				Size: int64(file.Buffer.Len()),
				Mode: 0766,
			})

			if err != nil {
				res.Errors++
				log.Println(err)
				return
			}

			_, err = io.Copy(tarbal, file.Buffer)

			if err != nil {
				res.Errors++
				log.Println(err)
			}
		}
	}()

	err = store.From.Walk(store.Loc, func(path string) {
		paths <- path
	})

	close(paths)
	read.Wait()
	close(files)
	write.Wait()
	tarbal.Close()
	gzip.Close()
	file.Close()

	if err != nil {
		return nil, err
	}

	export, err := store.To.Get(store.Dest)

	if err != nil {
		return nil, err
	}

	defer export.Close()

	if err := store.Remote.Put(store.Dest, export); err != nil {
		return nil, err
	}

	info, err := store.To.Stat(store.Dest)

	if err != nil {
		return nil, err
	}

	size := ((float64)(info.Size()) / 1024) / 1024

	switch req.ContentType {
	case pb.ContentType_JSON:
		proj.JSONSize = size
		proj.JSONPath = store.Dest
		proj.JSONAt = time.Now().UTC()
	case pb.ContentType_HTML:
		proj.HTMLSize = size
		proj.HTMLPath = store.Dest
		proj.HTMLAt = time.Now().UTC()
	case pb.ContentType_WIKITEXT:
		proj.WikitextSize = size
		proj.WikitextPath = store.Dest
		proj.WikitextAt = time.Now().UTC()
	}

	_, err = repo.Update(ctx, proj, func(q *orm.Query) *orm.Query {
		return q.WherePK()
	})

	return res, err
}
