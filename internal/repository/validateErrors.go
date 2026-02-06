package repository

import (
	"database/sql"

	"github.com/lib/pq"
)

var arrErr = map[pq.ErrorCode]error{
	"23505": ErrEnique,
	"23502": ErrNotNull,
	"53300": ErrMaxConn,
	"42501": ErrNoRights,
}

func ValidateErrorsPostgresql(err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return ValidationErrorsSql(err)
	}
	
	err, ok = arrErr[pqErr.Code]
	if !ok {
		return ErrDatabase
	}

	return err
}

func ValidationErrorsSql(err error) error {
	arr := map[error]error{
		sql.ErrNoRows: ErrNoFound,
	}

	value, ok := arr[err]
	if !ok {
		return ErrDatabase
	}

	return value
}
