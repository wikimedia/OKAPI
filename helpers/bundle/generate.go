package bundle

import (
	"bufio"
	"okapi/helpers/progress"
	"okapi/lib/bzip"
	"okapi/lib/storage"
	"okapi/models"
	"os"
	"time"
)

// Generate function to create bundle from existing file
func Generate(project *models.Project) error {
	path, err := bzip.Compress(project.BundlePath())

	if err != nil {
		return err
	}

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	reader := &progress.Reader{
		Message: "File '" + project.CompressedBundleName() + "' upload progress '%d%%'!",
		Reader:  bufio.NewReader(file),
		Size:    info.Size(),
	}

	if err = storage.Remote.Client().Put(project.RemoteBundlePath(), reader); err != nil {
		return err
	}

	project.Size = ((float64)(reader.Size / 1024)) / 1024
	project.Path = project.RemoteBundlePath()
	project.DumpedAt = time.Now()
	return models.Save(project)
}
