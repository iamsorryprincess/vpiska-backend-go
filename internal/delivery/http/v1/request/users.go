package request

import "regexp"

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

type CreateUserRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (request *CreateUserRequest) Validate() ([]string, error) {
	var validationErrors []string

	if request.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		return nil, err
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

	return validationErrors, nil
}

type LoginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (request *LoginUserRequest) Validate() ([]string, error) {
	var validationErrors []string

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	return validationErrors, nil
}

type ChangePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (request *ChangePasswordRequest) Validate() ([]string, error) {
	var validationErrors []string

	if request.ID == "" {
		validationErrors = append(validationErrors, emptyIDError)
	} else if matched, err := regexp.MatchString(idRegexp, request.ID); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidIdFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if request.Password != request.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors, nil
}
