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
	ErrUserBalanceLessThanSumToWithdraw  = errors.New("user balance less than sum to withdraw")
	ErrUserBalanceHasChanged             = errors.New("user balance has changed in other transaction")
	ErrOrderNumberOrSumIsNotFilled       = errors.New("order number or sum is not filled")
	ErrWithdrawalsWasNotFound            = errors.New("withdrawals to current user was not found")
)
