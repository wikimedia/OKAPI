package projects

import (
	"okapi/lib/env"
	"okapi/models"
	"os"
	"strings"
)

// GetHTMLPath get local html path
func GetHTMLPath(project *models.Project) string {
	return env.Context.VolumeMountPath + "/html/" + strings.Replace(project.SiteURL, "https://", "", 1)
}

// CreateExportFile create export file in file system
func CreateExportFile(project *models.Project) (*os.File, error) {
	for _, path := range []string{"/exports", "/exports/" + project.DBName} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(env.Context.VolumeMountPath+path, 0766)
		}
	}

	return os.Create(project.GetExportPath())
}
