package search

// Response api response for search
type Response struct {
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Total int         `json:"total"`
	Limit int         `json:"limit"`
	Items interface{} `json:"items"`
}
