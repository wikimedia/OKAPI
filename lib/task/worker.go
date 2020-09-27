package task

// Worker single unit processor
type Worker func(id int, payload Payload) (string, map[string]interface{}, error)
