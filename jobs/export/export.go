package export

import (
	"archive/tar"
	"okapi/helpers/bundle"
	"okapi/helpers/damaging"
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

	damaged, err := damaging.GetMap(ctx.Project.DBName)

	if err != nil {
		return nil, nil, nil, err
	}

	paths := make(chan *path)
	files := make(chan *file)

	gzip := pgzip.NewWriter(project)
	gzip.SetConcurrency(1<<20, runtime.NumCPU()*2)
	defer gzip.Close()

	tar := tar.NewWriter(gzip)
	defer tar.Close()

	wg := sync.WaitGroup{}
	path := projects.GetHTMLPath(ctx.Project)
	length := len(path) + 1

	for i := 0; i < ctx.Params.Workers; i++ {
		wg.Add(1)
		go readWorker(ctx, paths, files, &wg)
	}

	go writeWorker(ctx, files, tar, &wg)

	wg.Add(1)
	go walkPath(path, paths, &wg, length, damaged)

	wg.Wait()
	close(files)
	damaging.Delete(ctx.Project.DBName)

	return nil, nil, nil, bundle.Upload(ctx.Project)
}
