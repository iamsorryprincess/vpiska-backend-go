package v1

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorCode string `json:"errorCode"`
}

type ApiResponse struct {
	IsSuccess bool             `json:"isSuccess"`
	Result    interface{}      `json:"result"`
	Errors    []*ErrorResponse `json:"errors"`
}

func createDomainErrorResponse(err error) *ApiResponse {
	return &ApiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors: []*ErrorResponse{{
			ErrorCode: err.Error(),
		}},
	}
}

func createValidationErrorResponse(errs []string) *ApiResponse {
	errorsResponses := make([]*ErrorResponse, len(errs))

	for index, item := range errs {
		errorsResponses[index] = &ErrorResponse{
			ErrorCode: item,
		}
	}

	return &ApiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors:    errorsResponses,
	}
}

func createSuccessResponse(response interface{}) *ApiResponse {
	return &ApiResponse{
		IsSuccess: true,
		Result:    response,
		Errors:    nil,
	}
}

func writeDomainErrorResponse(writer http.ResponseWriter, domainError error) {
	errorResponse := createDomainErrorResponse(domainError)
	bytes, marshalErr := json.Marshal(errorResponse)

	if marshalErr != nil {
		panic(marshalErr)
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		panic(writeErr)
	}
}

func writeSuccessResponse(writer http.ResponseWriter, response interface{}) {
	bytes, marshalErr := json.Marshal(createSuccessResponse(response))

	if marshalErr != nil {
		panic(bytes)
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		panic(writeErr)
	}
}
