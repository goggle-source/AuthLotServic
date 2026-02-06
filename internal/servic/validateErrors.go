package servic

import "github.com/goggle-source/authLotServic/internal/repository"

var arrErr = map[error]error{
	repository.ErrDatabase: ErrInternal,
	repository.ErrEnique:   ErrClientEnuqie,
	repository.ErrMaxConn:  ErrInternal,
	repository.ErrNotNull:  ErrClientNotNull,
	repository.ErrPassword: ErrClientPassword,
	repository.ErrNoFound:  ErrNoFound,
}

func ValidationError(err error) error {
	resultErr, ok := arrErr[err]
	if !ok {
		return ErrInternal
	}

	return resultErr
}
