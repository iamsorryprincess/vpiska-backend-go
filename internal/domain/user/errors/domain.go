package errors

import "errors"

var (
	NameAlreadyUse         = errors.New("name already use")
	PhoneAlreadyUse        = errors.New("phone already use")
	NameAndPhoneAlreadyUse = errors.New("name and phone already use")
	UserNotFound           = errors.New("user not found")
	InvalidPassword        = errors.New("invalid password")
	InternalError          = errors.New("internal error")
)
