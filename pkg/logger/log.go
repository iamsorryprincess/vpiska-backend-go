package logger

import (
	"io"
	"log"
	"os"
	"runtime/debug"
)

type logger struct {
	logFile            *os.File
	errorLoggerConsole *log.Logger
	errorLoggerFile    *log.Logger
	infoLogger         *log.Logger
	warningLogger      *log.Logger
}

func NewLogLogger() (Logger, io.Closer, error) {
	_, err := os.Stat("logs")

	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.Mkdir("logs", 0777)
			if mkdirErr != nil {
				return nil, nil, mkdirErr
			}
		} else {
			return nil, nil, err
		}
	}

	logFile, err := os.OpenFile("logs/logs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, nil, err
	}

	return &logger{
		logFile:            logFile,
		errorLoggerFile:    log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLoggerConsole: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:         log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger:      log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, logFile, nil
}

func (l *logger) LogError(err error) {
	l.errorLoggerFile.Println(err, string(debug.Stack()))
	l.errorLoggerConsole.Println(err, string(debug.Stack()))
}

func (l *logger) LogErrorString(err string) {
	l.errorLoggerFile.Println(err, string(debug.Stack()))
}

func (l *logger) LogInfo(message string) {
	l.infoLogger.Println(message)
}

func (l *logger) LogWarning(message string) {
	l.warningLogger.Println(message)
}
