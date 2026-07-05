package domain

import "errors"

var (
	ErrEmail       = errors.New("such an email has already been registered")
	ErrUserNoFound = errors.New("user is not found")
	ErrPasswordOrEmail    = errors.New("invalid password or email")
)
