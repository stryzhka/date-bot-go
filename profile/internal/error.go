package profile

import "errors"

var (
	ErrValidationUserId      = errors.New("Invalid user id")
	ErrValidationName        = errors.New("Invalid user name")
	ErrUserNotFound          = errors.New("User not found")
	ErrUserAlreadyExists     = errors.New("User already exists")
	ErrValidationGender      = errors.New("Invalid gender")
	ErrValidationDescription = errors.New("Invalid description")
)
