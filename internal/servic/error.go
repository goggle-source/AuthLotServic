package servic

import "errors"

var (
	ErrInternal       = errors.New("servic error")
	ErrGenerateJWT    = errors.New("error when creating a token")
	ErrValidateToken  = errors.New("token validation error")
	ErrClientEnuqie   = errors.New("Such a user already exists")
	ErrClientNotNull  = errors.New("incorrect data entry")
	ErrClientPassword = errors.New("invalid password")
	ErrNoFound        = errors.New("user is not found")
)
