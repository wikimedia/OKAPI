package bundle

import (
	"bufio"
	"okapi/helpers/progress"
	"okapi/lib/storage"
	"okapi/models"
	"os"
	"time"
)

// Upload function to create bundle from existing file
func Upload(project *models.Project) error {
	path := project.GetExportPath()
	remotePath := project.GetRemoteExportPath()

	file, err := os.Open(path)

	if err != nil {
		return err
	}

	info, err := file.Stat()

	if err != nil {
		return err
	}

	reader := &progress.Reader{
		Message: "File '" + project.GetExportName() + "' upload progress '%d%%'!",
		Reader:  bufio.NewReader(file),
		Size:    info.Size(),
	}

	if err = storage.Remote.Client().Put(remotePath, reader); err != nil {
		return err
	}

	project.Size = ((float64)(reader.Size / 1024)) / 1024
	project.Path = remotePath
	project.DumpedAt = time.Now()

	_, err = models.DB().
		Model(project).
		Column("size", "path", "dumped_at", "updated_at").
		WherePK().
		Update()

	return err
}
