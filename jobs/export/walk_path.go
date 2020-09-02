package export

import "github.com/karrick/godirwalk"

func walkPath(path string, files chan string) error {
	return godirwalk.Walk(path, &godirwalk.Options{
		Unsorted: true,
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Callback: func(filePath string, de *godirwalk.Dirent) error {
			if !de.IsDir() {
				files <- filePath
			}

			return nil
		},
	})
}
