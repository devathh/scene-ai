package consts

import "errors"

var (
	// auth
	ErrInvalidFirstname = errors.New("firstname must be more than 3 and less than 64")
	ErrInvalidLastname  = errors.New("lastname must be more than 3 and less than 64")
	ErrInvalidPassword  = errors.New("password must be more than 6")
	ErrInvalidEmail     = errors.New("invalid email")
	// scenario
	ErrEmptyTitle             = errors.New("title cannot be empty")
	ErrEmptyScenarioPrompt    = errors.New("scenario prompt is empty")
	ErrEmptyGlobalStylePrompt = errors.New("global style prompt is empty")
	ErrEmptyVideoPrompt       = errors.New("video's prompt is empty")
	ErrNoScenes               = errors.New("scenario must be contains at least 1 scene")
	ErrInvalidID              = errors.New("invalid id")
	ErrInvalidStatus          = errors.New("status is invalid")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionNotFound   = errors.New("session not found")
	ErrScenarioNotFound  = errors.New("scenario not found")
	ErrInvalidLimit      = errors.New("limit must be more than 0 and less than 1000")

	ErrInvalidUserID = errors.New("invalid user's id")
	ErrInvalidToken  = errors.New("invalid token")
	ErrForbidden     = errors.New("access denied")

	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidIP          = errors.New("invalid ip")
	ErrInvalidUserAgent   = errors.New("invalid user-agent")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
