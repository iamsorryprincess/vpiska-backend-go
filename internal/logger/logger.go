package logger

import "log"

type Logger interface {
	LogError(err error)
}

type logger struct {
}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) LogError(err error) {
	log.Println(err)
}
