package filter

import (
	"github.com/go-pg/pg/v10/orm"
)

// Equal filter is used to query DB by the exact column value
func Equal(column string, param string) func(*orm.Query) {
	return func(query *orm.Query) {
		if len(param) > 0 {
			query.Where(column+" = ?", param)
		}
	}
}
