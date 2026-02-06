package grpc

import "errors"

var (
	ErrClientEnuqie   = errors.New("Such a user already exists")
	ErrClientNotNull  = errors.New("incorrect data entry")
	ErrClientPassword = errors.New("invalid password")
	ErrInternal       = errors.New("server error")
	ErrNoValidToken   = errors.New("invalid token")
	ErrRequired       = errors.New("argument is required")
	ErrNoCorrectEmail = errors.New("email is not correct")
	ErrArgumnetLength = errors.New("invalid length")
	ErrNoFound        = errors.New("user is not found")
)
