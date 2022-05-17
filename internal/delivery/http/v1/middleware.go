package v1

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
)

var unauthorizedResponse = createDomainErrorResponse(auth.ErrInvalidToken)

func (h *Handler) jwtAuth(context *gin.Context) {
	headerValue := context.GetHeader("Authorization")

	if headerValue == "" {
		context.AbortWithStatusJSON(http.StatusOK, unauthorizedResponse)
		return
	}

	encodedToken := strings.TrimPrefix(headerValue, "Bearer ")

	if encodedToken == "" {
		context.AbortWithStatusJSON(http.StatusOK, unauthorizedResponse)
		return
	}

	token, err := h.tokenManager.ParseToken(encodedToken)

	if err != nil {
		if err == auth.ErrInvalidToken {
			context.AbortWithStatusJSON(http.StatusOK, unauthorizedResponse)
			return
		}

		h.errorLogger.Println(err)
		context.AbortWithStatusJSON(http.StatusOK, createDomainErrorResponse(errInternal))
		return
	}

	validationErrs, err := validateId(token.ID)

	if err != nil {
		h.errorLogger.Println(err)
		context.AbortWithStatusJSON(http.StatusOK, createDomainErrorResponse(errInternal))
		return
	}

	if validationErrs != nil {
		context.AbortWithStatusJSON(http.StatusOK, createValidationErrorResponse(validationErrs))
		return
	}

	context.Set("UserID", token.ID)
	context.Set("Username", token.Name)
	context.Set("UserImage", token.ImageID)
}

func validateId(id string) ([]string, error) {
	var validationErrors []string

	if id == "" {
		validationErrors = append(validationErrors, emptyIDError)
	} else if matched, err := regexp.MatchString(idRegexp, id); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidIdFormatError)
	}

	return validationErrors, nil
}
