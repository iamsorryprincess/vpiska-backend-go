package logger

type Logger interface {
	LogError(err error)
	LogErrorString(err string)
	LogInfo(message string)
	LogWarning(message string)
}
