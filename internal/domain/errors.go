package domain

import "errors"

var (
	ErrPhoneAlreadyUse        = errors.New("PhoneAlreadyUse")
	ErrNameAlreadyUse         = errors.New("NameAlreadyUse")
	ErrNameAndPhoneAlreadyUse = errors.New("NameAndPhoneAlreadyUse")
	ErrUserNotFound           = errors.New("UserNotFound")
	ErrInvalidPassword        = errors.New("InvalidPassword")
	ErrInternalError          = errors.New("InternalError")
)

func MapUserError(err error) error {
	switch err {
	case ErrPhoneAlreadyUse, ErrNameAlreadyUse, ErrNameAndPhoneAlreadyUse, ErrUserNotFound, ErrInvalidPassword:
		return err
	default:
		return ErrInternalError
	}
}
