package v1

import "net/http"

type loggingWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func newLoggingWriter(writer http.ResponseWriter) *loggingWriter {
	return &loggingWriter{
		ResponseWriter: writer,
	}
}

func (w *loggingWriter) Write(body []byte) (int, error) {
	w.Body = body
	return w.ResponseWriter.Write(body)
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
