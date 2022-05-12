package domain

import "errors"

var (
	ErrPhoneAlreadyUse        = errors.New("PhoneAlreadyUse")
	ErrNameAlreadyUse         = errors.New("NameAlreadyUse")
	ErrNameAndPhoneAlreadyUse = errors.New("NameAndPhoneAlreadyUse")
	ErrUserNotFound           = errors.New("UserNotFound")
	ErrInvalidPassword        = errors.New("InvalidPassword")

	ErrEmptyMedia    = errors.New("MediaIsEmpty")
	ErrMediaNotFound = errors.New("MediaNotFound")

	ErrInternalError = errors.New("InternalError")
)

func MapDomainError(err error) error {
	switch err {
	case
		ErrPhoneAlreadyUse,
		ErrNameAlreadyUse,
		ErrNameAndPhoneAlreadyUse,
		ErrUserNotFound,
		ErrInvalidPassword,
		ErrEmptyMedia,
		ErrMediaNotFound:
		return err
	default:
		return ErrInternalError
	}
}
