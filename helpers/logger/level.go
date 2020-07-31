package logger

// Level alert level type
type Level uint8

// Alert levels
const (
	DEBUG         Level = 0
	INFO          Level = 1
	WARN          Level = 2
	ERROR         Level = 3
	FATAL         Level = 4
	UNKNOWN       Level = 5
	EMERGENCY     Level = 10
	ALERT         Level = 11
	CRITICAL      Level = 12
	WARNING       Level = 14
	NOTICE        Level = 15
	INFORMATIONAL Level = 16
)
