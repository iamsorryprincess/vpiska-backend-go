package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
)

type userID string

var userIdKey = userID("UserID")

var unauthorizedResponse = newErrorResponse(auth.ErrInvalidToken.Error())

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
