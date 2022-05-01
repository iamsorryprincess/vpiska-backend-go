package v1

import (
	"encoding/json"
	"net/http"
	"regexp"
)

const (
	phoneRegexp            = `^\d{10}\b$`
	requiredPasswordLength = 6
)

const (
	emptyIDError                = "IdIsEmpty"
	emptyNameError              = "NameIsEmpty"
	emptyPhoneError             = "PhoneIsEmpty"
	emptyPasswordError          = "PasswordIsEmpty"
	invalidIDFormatError        = "IdInvalidFormat"
	invalidPhoneFormatError     = "PhoneRegexInvalid"
	invalidPasswordLengthError  = "PasswordLengthInvalid"
	invalidConfirmPasswordError = "ConfirmPasswordInvalid"
)

type Validated interface {
	Validate() []string
}

func validateRequest(writer http.ResponseWriter, request Validated) bool {
	if validationErrors := request.Validate(); validationErrors != nil {
		errorResponse := createValidationErrorResponse(validationErrors)
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

		return false
	}

	return true
}

func (request *CreateUserRequest) Validate() []string {
	var validationErrors []string

	if request.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		panic(err)
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if request.Password != request.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors
}

func (request *LoginUserRequest) Validate() []string {
	var validationErrors []string

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		panic(err)
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	return validationErrors
}
