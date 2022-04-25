package errors

import "errors"

var (
	EmptyID                = errors.New("empty id")
	EmptyName              = errors.New("empty name")
	EmptyPhone             = errors.New("empty phone")
	EmptyPassword          = errors.New("empty password")
	InvalidPhoneRegexp     = errors.New("invalid phone format")
	InvalidPasswordLength  = errors.New("invalid password length")
	InvalidConfirmPassword = errors.New("invalid confirm password")
	InvalidIdFormat        = errors.New("invalid id format")
)
