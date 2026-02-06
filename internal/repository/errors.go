package repository

import "errors"

var (
	ErrEnique   = errors.New("duplicate key")
	ErrNotNull  = errors.New("the field is required")
	ErrDatabase = errors.New("database error")
	ErrMaxConn  = errors.New("the maximum number of connections is exceeded")
	ErrNoRights = errors.New("no rights")
	ErrPassword = errors.New("invalid password")
	ErrNoFound  = errors.New("user is not found")
)
