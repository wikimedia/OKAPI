package queue

// Worker queue processor worker
type Worker func(payload string) (string, map[string]interface{}, error)
