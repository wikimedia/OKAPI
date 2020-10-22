package success

// Success API success response
type Success struct {
	Message string `json:"message"`
}

// Message generate error message
func Message(message string) *Success {
	return &Success{
		Message: message,
	}
}
