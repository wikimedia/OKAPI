package exception

// Exception API error response
type Exception struct {
	Message string `json:"message"`
}

// Message generate error message
func Message(err error) *Exception {
	return &Exception{
		Message: err.Error(),
	}
}
