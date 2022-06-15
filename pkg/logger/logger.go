package logger

type Logger interface {
	LogError(err error)
	LogInfo(message string)
	LogWarning(message string)
	LogHttpRequest(url string, method string, body string, contentType string)
	LogHttpResponse(url string, method string, statusCode int, body string, contentType string)
}
