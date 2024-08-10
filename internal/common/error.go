package common

import "errors"

var (
	ErrUserDataIsNotValid                = errors.New("user data is not valid")
	ErrUserAlreadyExists                 = errors.New("user already exists")
	ErrOrdersWasNotFound                 = errors.New("orders to current user was not found")
	ErrOrderNumberIsNotFilled            = errors.New("order number is not filled")
	ErrOrderNumberIsNotValid             = errors.New("order number is not valid")
	ErrOrderNumberHasUploadedOtherUser   = errors.New("order number has uploaded from other user")
	ErrOrderNumberHasUploadedCurrentUser = errors.New("order number has uploaded early")
)
