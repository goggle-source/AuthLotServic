package servic

import "errors"

var (
	ErrInternal       = errors.New("servic error")
	ErrGenerateJWT    = errors.New("error when creating a token")
	ErrValidateToken  = errors.New("token validation error")
	ErrClientNotNull  = errors.New("incorrect data entry")
	ErrClientPassword = errors.New("invalid password")
	ErrLogin          = errors.New("invalid login")
)
