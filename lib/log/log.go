package log

// Log interface to send debug info
type Log interface {
	Info(shortMessage string, fullMessage string, params ...map[string]interface{})
	Success(shortMessage string, fullMessage string, params ...map[string]interface{})
	Error(shortMessage string, fullMessage string, params ...map[string]interface{})
	Panic(shortMessage string, fullMessage string, params ...map[string]interface{})
}
