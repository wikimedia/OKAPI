package logger

// Message log message structure
type Message struct {
	Version      string
	Host         string
	ShortMessage string
	FullMessage  string
	Job          string
	Level        Level
	Category     Category
	Params       map[string]interface{}
}
