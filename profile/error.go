package profile

import "errors"

var (
	ErrValidationUserId = errors.New("Invalid user id")
	ErrValidationName   = errors.New("Invalid user name")
)
