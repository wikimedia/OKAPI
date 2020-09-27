package queue

// Worker queue processor worker
type Worker func(payload string) (string, map[string]interface{}, error)

func runWorker(ctx *Context, list chan string, worker Worker) {
	for item := range list {
		message, info, err := worker(item)

		if err != nil {
			ctx.Log.Error("queue worker error", err.Error(), info)
		} else {
			ctx.Log.Info("queue worker success", message, info)
		}
	}
}
