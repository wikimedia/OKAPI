package queue

import "okapi/helpers/logger"

// Worker queue processor worker
type Worker func(payload string) (string, map[string]interface{}, error)

func runWorker(list chan string, worker Worker) {
	for item := range list {
		message, info, err := worker(item)
		if err != nil {
			logger.Queue.Error("Queue worker encountered and error!", err.Error(), info)
		} else {
			logger.Queue.Success("Queue worker processed the unit!", message, info)
		}
	}
}
