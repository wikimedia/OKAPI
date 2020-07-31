package search

// Order order param
type Order string

// ASC standart ordering
const ASC Order = "asc"

// DESC reverse ordering
const DESC Order = "desc"

// Request to api struct
type Request struct {
	Page  int    `json:"page" form:"page"`
	Limit int    `json:"limit" form:"limit"`
	Sort  string `json:"sort" form:"sort"`
	Order Order  `json:"order" form:"order"`
}
