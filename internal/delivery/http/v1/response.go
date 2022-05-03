package v1

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
