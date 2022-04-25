package v1

type ApiResponse[T any] struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
	Result    T      `json:"result"`
}

func createSuccessResponse[T any](result T) *ApiResponse[T] {
	return &ApiResponse[T]{
		IsSuccess: true,
		Result:    result,
	}
}

func createErrorResponse[T any](err error) *ApiResponse[T] {
	return &ApiResponse[T]{
		IsSuccess: false,
		Error:     err.Error(),
	}
}
