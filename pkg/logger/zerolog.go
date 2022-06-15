package logger

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
)

type zeroLogger struct {
	logger zerolog.Logger
}

func NewZeroLogger() (Logger, *os.File, error) {
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

	logFile, err := os.OpenFile("logs/logs.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, nil, err
	}

	return &zeroLogger{
		logger: zerolog.New(logFile).With().Timestamp().Logger(),
	}, logFile, nil
}

func (l *zeroLogger) LogInfo(message string) {
	l.logger.Info().Msg(message)
}

func (l *zeroLogger) LogWarning(message string) {
	l.logger.Warn().Msg(message)
}

func (l *zeroLogger) LogError(err error) {
	l.logger.Error().Err(err).Str("stacktrace", string(debug.Stack())).Msg("")
}

func (l *zeroLogger) LogHttpRequest(url string, method string, body string, contentType string) {
	l.logger.Trace().
		Str("type", "request").
		Str("url", url).
		Str("method", method).
		Str("body", body).
		Str("contentType", contentType).
		Msg("")
}

func (l *zeroLogger) LogHttpResponse(url string, method string, statusCode int, body string, contentType string) {
	l.logger.Trace().
		Str("type", "response").
		Str("url", url).
		Str("method", method).
		Str("body", body).
		Str("contentType", contentType).
		Int("status", statusCode).
		Msg("")
}
