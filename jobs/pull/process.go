package pull

import (
	"okapi/models"
	"sync"
)

func processor(page *models.Page, get getter, errs *errors, wg *sync.WaitGroup) {
	defer wg.Done()
	err := get(page)
	errs.Lock()
	defer errs.Unlock()
	errs.items = append(errs.items, err)
}
