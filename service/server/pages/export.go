package pages

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
	"okapi-data-service/models"
	"okapi-data-service/schema/v3"
	pb "okapi-data-service/server/pages/protos"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/protsack-stephan/dev-toolkit/pkg/repository"
	"github.com/protsack-stephan/dev-toolkit/pkg/storage"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/klauspost/pgzip"
)

// ExportMaxFileSize max file size for single export ndjson file
const ExportMaxFileSize = 1000000000 * 10

// ErrExportFileIsNil signal that export returned nil (for stubs and new interfaces)
var ErrExportFileIsNil = errors.New("io.ReadWriteCloser is equal to nil")

// ErrWrongTmpFileType error that shows that tmp file for exports is of wrong type
var ErrWrongTmpFileType = errors.New("tmp file is not *os.File of type")

// ExportStorage storage to get data from and to put archive into
type ExportStorage struct {
	Loc      string // path to the files directory (where to get data from)
	Dest     string // path to the destination directory (where to put the dump)
	MetaDest string // path to the meta destination directory (where to put the dump meta)
	To       interface {
		storage.Getter
		storage.Stater
		storage.Creator
		storage.Deleter
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

type ndFile struct {
	name string
	path string
}

type tmpFile struct {
	file *os.File
	size int64
}

func newTmpFile(rc io.ReadCloser) (*tmpFile, error) {
	tmpf := new(tmpFile)

	if file, ok := rc.(*os.File); ok {
		tmpf.file = file
	}

	if tmpf.file == nil {
		return nil, ErrWrongTmpFileType
	}

	return tmpf, nil
}

// Export generate new export file and upload it to the storage
func Export(ctx context.Context, req *pb.ExportRequest, repo exportRepo, store *ExportStorage) (*pb.ExportResponse, error) {
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
	writeWg, readWg, copyWg := new(sync.WaitGroup), new(sync.WaitGroup), new(sync.WaitGroup)
	files := make(chan []byte, int(req.Workers))
	paths := make(chan string, int(req.Workers))
	ndfiles := make(chan ndFile, int(req.Workers))
	failed := new(sync.Map)
	gzip := pgzip.NewWriter(file)

	if err := gzip.SetConcurrency(1<<20, runtime.NumCPU()*2); err != nil {
		return nil, err
	}

	tarbal := tar.NewWriter(gzip)

	// spin up workers that will be responsible for reading files from the filesystem
	readWg.Add(int(req.Workers))
	for i := 1; i <= int(req.Workers); i++ {
		go func() {
			defer readWg.Done()
			for path := range paths {
				file, err := store.From.Get(path)

				if err != nil {
					log.Println(err)
					continue
				}

				data, err := ioutil.ReadAll(file)
				_ = file.Close()

				if err != nil {
					log.Printf("path: %s, err: %v", path, err)
					continue
				}

				page := new(schema.Page)

				if err := json.Unmarshal(data, page); err != nil {
					log.Printf("path: %s, err: %v", path, err)
					failed.Store(path[len(store.Loc)+1:len(path)-len(".json")], struct{}{})
					continue
				}

				if page.Namespace.Identifier == int(req.Ns) {
					files <- data
				}
			}
		}()
	}

	// create sepperate worker that will get ndjson file and copy them to the tar.gz file
	copyWg.Add(1)
	go func() {
		defer copyWg.Done()

		for file := range ndfiles {
			copyf, err := store.To.Get(file.path)

			if err != nil {
				log.Println(err)
				continue
			}

			finfo, err := store.To.Stat(file.path)

			if err != nil {
				log.Println(err)
				continue
			}

			header := &tar.Header{
				Name:    file.name,
				Size:    finfo.Size(),
				Mode:    0766,
				ModTime: time.Now().UTC(),
			}

			if err := tarbal.WriteHeader(header); err != nil {
				log.Println(err)
			} else if _, err = io.Copy(tarbal, copyf); err != nil {
				log.Println(err)
			}

			if err := store.To.Delete(file.path); err != nil {
				log.Println(err)
			}
		}
	}()

	// create write worker
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()

		var tmpf *tmpFile
		headf, numf := 0, 0

		copy := func(numf int) {
			headf++
			name := fmt.Sprintf("%s_%d.ndjson", req.DbName, numf-1)
			ndfiles <- ndFile{name, fmt.Sprintf("tmp/%d/%s", req.Ns, name)}
		}

		for file := range files {
			if tmpf == nil || tmpf.size >= ExportMaxFileSize {
				if tmpf != nil {
					copy(numf)
				}

				storef, err := store.To.Create(fmt.Sprintf("tmp/%d/%s_%d.ndjson", req.Ns, req.DbName, numf))

				if err != nil {
					log.Println(err)
					continue
				}

				tmpf, err = newTmpFile(storef)

				if err != nil {
					log.Println(err)
					continue
				}

				numf++
			}

			res.Total++

			if err != nil {
				res.Errors++
				log.Println(err)
				continue
			}

			if _, err := tmpf.file.Write(file); err != nil {
				res.Errors++
				log.Println(err)
				continue
			}

			tmpf.size += int64(len(file))

			nl := []byte("\n")

			if _, err := tmpf.file.Write(nl); err != nil {
				res.Errors++
				log.Println(err)
				continue
			}

			tmpf.size += int64(len(nl))
		}

		if numf != headf {
			copy(numf)
		}
	}()

	err = store.From.Walk(store.Loc, func(path string) {
		paths <- path
	})

	close(paths)
	readWg.Wait()
	close(files)
	writeWg.Wait()
	close(ndfiles)
	copyWg.Wait()
	_ = tarbal.Close()
	_ = gzip.Close()
	_ = file.Close()

	if err != nil {
		return nil, err
	}

	if res.Total == 0 {
		return res, nil
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

	// generate body md5 hash
	h := md5.New() // #nosec G401
	if _, err = io.Copy(h, export); err != nil {
		return nil, err
	}

	size := math.Round((((float64)(info.Size())/1024)/1024)*100) / 100
	version := fmt.Sprintf("%x", h.Sum(nil))
	datetime := time.Now().UTC()
	meta := schema.Project{
		Name:         proj.SiteName,
		Identifier:   proj.DbName,
		URL:          proj.SiteURL,
		Version:      &version,
		DateModified: &datetime,
		Size: &schema.Size{
			Value:    size,
			UnitText: "MB",
		},
	}

	if proj.Language != nil {
		meta.InLanguage = &schema.Language{
			Name:       proj.Language.LocalName,
			Identifier: proj.Language.Code,
		}
	}

	metadata, err := json.Marshal(meta)

	if err != nil {
		return nil, err
	}

	if err := store.Remote.Put(store.MetaDest, bytes.NewReader(metadata)); err != nil {
		return nil, err
	}

	titles := []string{}

	failed.Range(func(key, value interface{}) bool {
		switch title := key.(type) {
		case string:
			titles = append(titles, title)
		}

		return true
	})

	if len(titles) > 0 {
		query := func(q *orm.Query) *orm.Query {
			return q.Set("failed = true").Where("db_name = ? and title in (?)", req.DbName, pg.Strings(titles))
		}

		if _, err := repo.Update(ctx, new(models.Page), query); err != nil {
			log.Printf("titles: %v, err: %v", titles, err)
		}
	}

	return res, nil
}
