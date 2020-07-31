package bzip

import (
	"os"
	"os/exec"
)

// Compress bzip compress function
func Compress(filePath string) (string, error) {
	// Remove file if it exists
	os.Remove(filePath + ".bz2")

	cmd := exec.Command("bzip2", "-k", filePath)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return filePath + ".bz2", nil
}
