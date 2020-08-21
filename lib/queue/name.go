package queue

import "time"

// Name queue name type
type Name string

// Queue names
const (
	PageDelete   Name = "page_delete"
	PagePull     Name = "page_pull"
	PageRevision Name = "page_revision"
)

// Add push to the queue by name
func (name Name) Add(values ...interface{}) {
	Add(string(name), values)
}

// Pop get items from queue by name
func (name Name) Pop(frequency time.Duration) []string {
	return Pop(string(name), frequency)
}
