package export

import (
	"archive/tar"
	"okapi/helpers/bundle"
	"okapi/helpers/projects"
	"okapi/lib/task"
	"runtime"
	"sync"

	"github.com/klauspost/pgzip"
)

// Name task name for trigger
var Name task.Name = "export"

// Task for bundling the html
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	project, err := projects.CreateExportFile(ctx.Project)

	if err != nil {
		return nil, nil, nil, err
	}

	paths := make(chan string)
	files := make(chan *file)

	gzip := pgzip.NewWriter(project)
	gzip.SetConcurrency(1<<20, runtime.NumCPU()*2)
	defer gzip.Close()

	tar := tar.NewWriter(gzip)
	defer tar.Close()

	group := sync.WaitGroup{}
	path := projects.GetHTMLPath(ctx.Project)
	length := len(path) + 1

	for i := 0; i < *ctx.Cmd.Workers; i++ {
		go readWorker(paths, length, files, &group)
	}

	go writeWorker(files, tar, &group)

	err = walkPath(path, paths)

	close(paths)
	group.Wait()
	close(files)

	if err == nil {
		err = bundle.Upload(ctx.Project)
	}

	return nil, nil, nil, err
}
