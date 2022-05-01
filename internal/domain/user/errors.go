package user

import "errors"

var (
	PhoneAlreadyUse        = errors.New("PhoneAlreadyUse")
	NameAlreadyUse         = errors.New("NameAlreadyUse")
	NameAndPhoneAlreadyUse = errors.New("NameAndPhoneAlreadyUse")
	NotFound               = errors.New("UserNotFound")
	InvalidPassword        = errors.New("InvalidPassword")
	InternalError          = errors.New("InternalError")
)

func MapError(err error) error {
	switch err {
	case PhoneAlreadyUse, NameAlreadyUse, NameAndPhoneAlreadyUse, NotFound, InvalidPassword:
		return err
	default:
		return InternalError
	}
}
