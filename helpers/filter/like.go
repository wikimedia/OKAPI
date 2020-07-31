package filter

import (
	"github.com/go-pg/pg/v9/orm"
)

// Like filter is used to query DB by text chunk
func Like(column string, param string) func(*orm.Query) {
	return func(query *orm.Query) {
		if len(param) > 0 {
			query.Where(column+" ILIKE ?", "%"+param+"%")
		}
	}
}
