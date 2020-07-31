package queue

import "time"

// Name queue name type
type Name string

// Queue names
const (
	Sync       Name = "sync"
	Scan       Name = "scan"
	DeletePage Name = "delete-page"
	Test       Name = "test"
)

// Add push to the queue by name
func (name Name) Add(values ...interface{}) {
	Add(string(name), values)
}

// Pop get items from queue by name
func (name Name) Pop(frequency time.Duration) []string {
	return Pop(string(name), frequency)
}
