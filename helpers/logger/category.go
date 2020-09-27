package logger

// Type of log category
const (
	System  Category = "system"
	Job     Category = "job"
	API     Category = "api"
	Steream Category = "stream"
	Queue   Category = "queue"
	Runner  Category = "runner"
	Search  Category = "search"
)

// Category log category level
type Category string

// Panic method to panic for predefined category
func (category Category) Panic(shortMessage string, fullMessage string, params ...map[string]interface{}) {
	Panic(Message{
		Category:     category,
		ShortMessage: shortMessage,
		FullMessage:  fullMessage,
		Params:       getParams(params),
	})
}

// Error function to show and log the error by caregory
func (category Category) Error(shortMessage string, fullMessage string, params ...map[string]interface{}) {
	Error(Message{
		Category:     category,
		ShortMessage: shortMessage,
		FullMessage:  fullMessage,
		Params:       getParams(params),
	})
}

// Success function to show and log success message
func (category Category) Success(shortMessage string, fullMessage string, params ...map[string]interface{}) {
	Success(Message{
		Category:     category,
		ShortMessage: shortMessage,
		FullMessage:  fullMessage,
		Params:       getParams(params),
	})
}

// Send function to send the request message
func (category Category) Send(shortMessage string, fullMessage string, params ...map[string]interface{}) {
	Send(Message{
		Category:     category,
		ShortMessage: shortMessage,
		FullMessage:  fullMessage,
		Params:       getParams(params),
	})
}

// Info print info message to console
func (category Category) Info(shortMessage string, fullMessage string, params ...map[string]interface{}) {
	Info(Message{
		Category:     category,
		ShortMessage: shortMessage,
		FullMessage:  fullMessage,
		Params:       getParams(params),
	})
}

func getParams(params []map[string]interface{}) map[string]interface{} {
	for _, param := range params {
		return param
	}
	return map[string]interface{}{}
}

// Log simple function to go error with message
func (category Category) Log(err error, shortMessage string) {
	Error(Message{
		Category:     category,
		FullMessage:  err.Error(),
		ShortMessage: shortMessage,
	})
}
