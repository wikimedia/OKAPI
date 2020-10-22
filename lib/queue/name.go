package queue

import "time"

// Name queue name type
type Name string

// Queue names
const (
	PageDelete   Name = "page_delete"
	PagePull     Name = "page_pull"
	PageRevision Name = "page_revision"
	PageScore    Name = "page_score"
)

// Add push to the queue by name
func (name Name) Add(values ...interface{}) {
	Add(string(name), values)
}

// Pop get items from queue by name
func (name Name) Pop(frequency time.Duration) []string {
	return Pop(string(name), frequency)
}

// Subscribe set items to process
func (name Name) Subscribe(ctx *Context, worker Worker) {
	list := make(chan string)
	defer close(list)

	for i := 0; i < ctx.Workers; i++ {
		go runWorker(ctx, list, worker)
	}

	for {
		items := name.Pop(1 * time.Second)

		for _, item := range items {
			if item != getName(string(name)) {
				list <- item
			}
		}
	}
}
