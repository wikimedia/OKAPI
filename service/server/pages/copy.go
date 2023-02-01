package pages

import (
	"context"
	"fmt"
	pb "okapi-data-service/server/pages/protos"

	"github.com/protsack-stephan/dev-toolkit/pkg/storage"
)

// copyNumWorkers is the default number of workers for concurrency.
const copyNumWorkers = 10

// Copy copies project dump and metadata, as well as global exports metadata for group consumption.
// e.g., export/enwiki/enwiki_14.json -> export/enwiki/enwiki_group_1_14.json
// export/enwiki/enwiki_json_0.tar.gz -> export/enwiki/enwiki_group_1_json_0.tar.gz
// public/exports_0.json -> public/exports_group_1_0.json
func Copy(ctx context.Context, req *pb.CopyRequest, store storage.CopierWithContext, suffix string) (*pb.CopyResponse, error) {
	if req.Workers == 0 {
		req.Workers = copyNumWorkers
	}

	res := new(pb.CopyResponse)
	paths := make(map[string]string) // Store source-destination path pairs to copy

	for _, db := range req.DbNames {
		paths[fmt.Sprintf("export/%s/%s_%d.json", db, db, req.Ns)] = fmt.Sprintf("export/%s/%s%s_%d.json", db, db, suffix, req.Ns)
		paths[fmt.Sprintf("export/%s/%s_json_%d.tar.gz", db, db, req.Ns)] = fmt.Sprintf("export/%s/%s%s_json_%d.tar.gz", db, db, suffix, req.Ns)
	}

	size := len(paths)
	res.Total = int32(size)

	srcs := make(chan string, size)
	errs := make(chan error, size)

	for src := range paths {
		srcs <- src
	}

	close(srcs)

	// Read from paths channel and copy out each file, pushing the error to error channel.
	for i := 0; i < int(req.Workers); i++ {
		go func() {
			for src := range srcs {
				errs <- store.CopyWithContext(ctx, src, paths[src])
			}
		}()
	}

	for i := 0; i < size; i++ {
		if err := <-errs; err != nil {
			res.Errors++
		}
	}

	close(errs)

	return res, nil
}
