package queue

import (
	"time"
)

// Subscribe func to create queue subscriber (go routine use recommended)
func Subscribe(subscriber Name, worker Worker, workers int) {
	list := make(chan string)
	defer close(list)

	for i := 0; i < workers; i++ {
		go runWorker(list, worker)
	}

	for {
		items := subscriber.Pop(1 * time.Second)
		for _, item := range items {
			if item != getName(string(subscriber)) {
				list <- item
			}
		}
	}
}
