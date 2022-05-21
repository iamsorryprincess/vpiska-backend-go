package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

var errInternal = errors.New("InternalError")

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
}

type apiResponse struct {
	IsSuccess bool            `json:"isSuccess"`
	Errors    []errorResponse `json:"errors"`
	Result    interface{}     `json:"result"`
}

func createDomainErrorResponse(err error) apiResponse {
	return apiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors: []errorResponse{{
			ErrorCode: err.Error(),
		}},
	}
}

func createValidationErrorResponse(errs []string) apiResponse {
	errorsResponses := make([]errorResponse, len(errs))

	for index, item := range errs {
		errorsResponses[index] = errorResponse{
			ErrorCode: item,
		}
	}

	return apiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors:    errorsResponses,
	}
}

func createSuccessResponse(response interface{}) apiResponse {
	return apiResponse{
		IsSuccess: true,
		Result:    response,
		Errors:    nil,
	}
}

func writeErrorResponse(err error, logger logger.Logger, context *gin.Context) {
	var response apiResponse

	if domain.IsInternalError(err) {
		logger.LogError(err)
		response = createDomainErrorResponse(errInternal)
	} else {
		response = createDomainErrorResponse(err)
	}

	context.JSON(http.StatusOK, response)
}

func writeValidationErrResponse(validationErrs []string, context *gin.Context) {
	response := createValidationErrorResponse(validationErrs)
	context.JSON(http.StatusOK, response)
}

func writeResponse(result interface{}, context *gin.Context) {
	response := createSuccessResponse(result)
	context.JSON(http.StatusOK, response)
}
