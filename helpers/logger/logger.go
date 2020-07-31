package logger

import (
	"log"

	"gopkg.in/gookit/color.v1"
)

// Success function to print and log success info
func Success(message Message) {
	print(message, color.Success.Println)
	message.Level = INFO
	Send(message)
}

// Info function to print and log success info
func Info(message Message) {
	print(message, color.Info.Println)
}

// Error function to show and log an error message
func Error(message Message) {
	print(message, color.Error.Println)
	message.Level = WARN
	Send(message)
}

// Panic function to stop the execution
func Panic(message Message) {
	print(message, color.Error.Println)
	message.Level = ERROR
	Send(message)
	log.Panic(message.ShortMessage)
}

func print(message Message, printer func(a ...interface{})) {
	printer(message.ShortMessage)
	if len(message.FullMessage) > 0 {
		printer(message.FullMessage)
	}
}
