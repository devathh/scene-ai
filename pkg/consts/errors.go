package consts

import "errors"

var (
	ErrInvalidFirstname = errors.New("firstname must be more than 3 and less than 64")
	ErrInvalidLastname  = errors.New("lastname must be more than 3 and less than 64")
	ErrInvalidPassword  = errors.New("password must be more than 6")
	ErrInvalidEmail     = errors.New("invalid email")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionNotFound   = errors.New("session not found")

	ErrInvalidUserID = errors.New("invalid user's id")
	ErrInvalidToken  = errors.New("invalid token")

	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidIP          = errors.New("invalid ip")
	ErrInvalidUserAgent   = errors.New("invalid user-agent")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
