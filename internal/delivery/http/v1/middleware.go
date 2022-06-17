package v1

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
)

const invalidMethod = "invalid method"

type userID string

var userIdKey = userID("UserID")
var unauthorizedResponse = newErrorResponse(auth.ErrInvalidToken.Error())

func (h *Handler) GET(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			h.writeJSONResponse(writer, newErrorResponse(invalidMethod))
			return
		}
		next.ServeHTTP(writer, request)
	}
}

func (h *Handler) POST(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			h.writeJSONResponse(writer, newErrorResponse(invalidMethod))
			return
		}
		next.ServeHTTP(writer, request)
	}
}

func (h *Handler) DELETE(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodDelete {
			h.writeJSONResponse(writer, newErrorResponse(invalidMethod))
			return
		}
		next.ServeHTTP(writer, request)
	}
}

func (h *Handler) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				var err error

				switch t := recovered.(type) {
				case error:
					err = t
					break
				case string:
					err = errors.New(t)
					break
				default:
					err = errors.New("unknown error")
					break
				}

				if err != nil {
					h.logger.LogError(err)
				}

				h.writeJSONResponse(writer, newErrorResponse(internalError))
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func (h *Handler) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.RequestURI, "/api/") && !strings.Contains(request.RequestURI, "/websockets/") {
			var buf bytes.Buffer
			tee := io.TeeReader(request.Body, &buf)
			body, err := ioutil.ReadAll(tee)
			if err != nil {
				h.logger.LogError(err)
				h.writeJSONResponse(writer, newErrorResponse(internalError))
				return
			}

			requestContentType := request.Header.Get("Content-Type")
			switch requestContentType {
			case contentTypeJSON:
				h.logger.LogHttpRequest(request.RequestURI, request.Method, string(body), requestContentType)
				break
			default:
				h.logger.LogHttpRequest(request.RequestURI, request.Method, "(hidden)", requestContentType)
				break
			}

			request.Body = io.NopCloser(&buf)
			writerWithState := newLoggingWriter(writer)
			next.ServeHTTP(writerWithState, request)
			responseContentType := writer.Header().Get("Content-Type")
			switch responseContentType {
			case contentTypeJSON:
				h.logger.LogHttpResponse(request.RequestURI, request.Method, writerWithState.StatusCode, string(writerWithState.Body), responseContentType)
				break
			default:
				h.logger.LogHttpResponse(request.RequestURI, request.Method, writerWithState.StatusCode, "(hidden)", responseContentType)
				break
			}
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}

func (h *Handler) jwtAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		headerValue := request.Header.Get("Authorization")

		if headerValue == "" {
			h.writeJSONResponse(writer, unauthorizedResponse)
			return
		}

		encodedToken := strings.TrimPrefix(headerValue, "Bearer ")

		if encodedToken == "" {
			h.writeJSONResponse(writer, unauthorizedResponse)
			return
		}

		token, err := h.tokenManager.ParseToken(encodedToken)

		if err != nil {
			if err == auth.ErrInvalidToken {
				h.writeJSONResponse(writer, unauthorizedResponse)
				return
			}

			h.logger.LogError(err)
			h.writeJSONResponse(writer, newErrorResponse(internalError))
			return
		}

		validationErrs, err := validateId(token.ID)

		if err != nil {
			h.logger.LogError(err)
			h.writeJSONResponse(writer, newErrorResponse(internalError))
			return
		}

		if len(validationErrs) > 0 {
			h.writeJSONResponse(writer, newValidationErrsResponse(validationErrs))
			return
		}

		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), userIdKey, userID(token.ID))))
	}
}

func getUserID(request *http.Request) string {
	value, ok := request.Context().Value(userIdKey).(userID)

	if !ok {
		return ""
	}

	return string(value)
}
