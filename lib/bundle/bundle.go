package bundle

import (
	"os"
	"os/exec"
	"path/filepath"
)

// UpdateFileInTar updates file in tar archive
func UpdateFileInTar(tarPath string, filePath string) error {
	fileName := filepath.Base(filePath)

	// Copy and remove file
	// Workaround to flatten the archive structure
	// --transform configuration doesn't work because of the bug in tar library
	copyFile := exec.Command("cp", filePath, "./")
	archive := exec.Command("tar", "uvf", tarPath, fileName)
	removeCopy := exec.Command("rm", fileName)

	// Attach error logs use copyFile.Stdout = os.Stdout or similar to attach other logs
	copyFile.Stderr = os.Stdout
	archive.Stderr = os.Stdout
	removeCopy.Stderr = os.Stdout

	if err := copyFile.Run(); err != nil {
		return err
	}

	if err := archive.Run(); err != nil {
		return err
	}

	if err := removeCopy.Run(); err != nil {
		return err
	}

	return nil
}
