package v1

import "errors"

const (
	idRegexp               = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	phoneRegexp            = `^\d{10}\b$`
	requiredPasswordLength = 6
)

const (
	emptyIDError                = "IdIsEmpty"
	emptyNameError              = "NameIsEmpty"
	emptyPhoneError             = "PhoneIsEmpty"
	emptyPasswordError          = "PasswordIsEmpty"
	invalidIdFormatError        = "InvalidIdFormat"
	invalidPhoneFormatError     = "PhoneRegexInvalid"
	invalidPasswordLengthError  = "PasswordLengthInvalid"
	invalidConfirmPasswordError = "ConfirmPasswordInvalid"
)

var errEmptyId = errors.New("IdIsEmpty")

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
}

type apiResponse struct {
	IsSuccess bool            `json:"isSuccess"`
	Result    interface{}     `json:"result"`
	Errors    []errorResponse `json:"errors"`
}

type requestForValidate interface {
	validate() ([]string, error)
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
