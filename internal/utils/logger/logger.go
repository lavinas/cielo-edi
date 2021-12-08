package logger

import (
	"log"
)


type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l Logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l Logger) Println(v ...interface{}) {
	log.Println(v...)
}
