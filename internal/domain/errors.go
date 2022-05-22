package domain

import "errors"

var (
	ErrPhoneAlreadyUse        = errors.New("PhoneAlreadyUse")
	ErrNameAlreadyUse         = errors.New("NameAlreadyUse")
	ErrNameAndPhoneAlreadyUse = errors.New("NameAndPhoneAlreadyUse")
	ErrUserNotFound           = errors.New("UserNotFound")
	ErrInvalidPassword        = errors.New("InvalidPassword")

	ErrMediaNotFound = errors.New("MediaNotFound")

	ErrEventNotFound = errors.New("EventNotFound")
)

func IsInternalError(err error) bool {
	switch err {
	case
		ErrPhoneAlreadyUse,
		ErrNameAlreadyUse,
		ErrNameAndPhoneAlreadyUse,
		ErrUserNotFound,
		ErrInvalidPassword,
		ErrMediaNotFound,
		ErrEventNotFound:
		return false
	default:
		return true
	}
}
