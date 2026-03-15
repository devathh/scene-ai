package consts

import "errors"

var (
	ErrInvalidFirstname = errors.New("firstname must be more than 3 and less than 64")
	ErrInvalidLastname  = errors.New("lastname must be more than 3 and less than 64")
	ErrInvalidPassword  = errors.New("password must be more than 6")

	ErrUserNotFound = errors.New("user not found")
)
