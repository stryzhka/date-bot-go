package matching

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrAlreadyLiked = errors.New("Already liked")
	ErrLikeNotFound = errors.New("Like not found")
	ErrAutoLike     = errors.New("Can't like yourself")
)
