package v1

import (
	"errors"
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

func newErrorResponse(err string) apiResponse {
	return apiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors: []errorResponse{{
			ErrorCode: err,
		}},
	}
}

func newValidationErrsResponse(errs []string) apiResponse {
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

func newSuccessResponse(response interface{}) apiResponse {
	return apiResponse{
		IsSuccess: true,
		Result:    response,
		Errors:    nil,
	}
}
