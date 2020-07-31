package logger

// Type of log category
const (
	SYSTEM Category = "system"
	JOB    Category = "job"
	API    Category = "api"
	STREAM Category = "stream"
	QUEUE  Category = "queue"
	RUNNER Category = "runner"
)

// Category log category level
type Category string

// Panic method to panic for predefined category
func (category Category) Panic(message Message) {
	message.Category = category
	Panic(message)
}

// Error function to show and log the error by caregory
func (category Category) Error(message Message) {
	message.Category = category
	Error(message)
}

// Success function to show and log success message
func (category Category) Success(message Message) {
	message.Category = category
	Success(message)
}

// Send function to send the request message
func (category Category) Send(message Message) {
	message.Category = category
	Send(message)
}

// Info print info message to console
func (category Category) Info(message Message) {
	message.Category = category
	Info(message)
}
