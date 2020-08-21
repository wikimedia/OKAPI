package bundle

import (
	"okapi/models"
	"os"
	"os/exec"
)

// Delete file from tar
func Delete(project *models.Project, page *models.Page) error {
	delete := exec.Command("tar", "--delete", "--file="+project.BundlePath(), page.Title+".html")
	delete.Stderr = os.Stdout
	return delete.Run()
}
