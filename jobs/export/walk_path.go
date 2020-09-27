package export

import (
	"sync"

	"github.com/karrick/godirwalk"
)

func walkPath(pathname string, paths chan *path, wg *sync.WaitGroup, length int, damaged map[string]bool) {
	callback := func(filePath string, de *godirwalk.Dirent) error {
		if !de.IsDir() {
			paths <- &path{
				file: filePath[length:],
				full: filePath,
			}
		}

		return nil
	}

	if len(damaged) > 0 {
		callback = func(filePath string, de *godirwalk.Dirent) error {
			if !de.IsDir() {
				title := filePath[length:]

				if _, ok := damaged[title[:len(title)-5]]; !ok {
					paths <- &path{
						file: title,
						full: filePath,
					}
				}
			}

			return nil
		}
	}

	godirwalk.Walk(pathname, &godirwalk.Options{
		Unsorted: true,
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Callback: callback,
	})

	close(paths)
	wg.Done()
}
