package queue

import (
	"time"

	"okapi/helpers/logger"
	"okapi/lib/cmd"
)

// Subscribe func to create queue subscriber (go routine use recommended)
func Subscribe(subscriber Name, worker Worker) {
	list := make(chan string)
	defer close(list)

	for i := 0; i < *cmd.Context.Workers; i++ {
		go func() {
			for item := range list {
				message, info, err := worker(item)
				if err != nil {
					logger.QUEUE.Error(logger.Message{
						ShortMessage: "Queue worker encountered and error!",
						FullMessage:  err.Error(),
						Params:       info,
					})
				} else {
					logger.QUEUE.Success(logger.Message{
						ShortMessage: "Queue worker processed the unit!",
						FullMessage:  message,
						Params:       info,
					})
				}
			}
		}()
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
