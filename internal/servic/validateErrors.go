package servic

import "github.com/goggle-source/authLotServic/internal/repository"

var arrErr = map[error]error{
	repository.ErrDatabase: ErrInternal,
	repository.ErrEnique:   ErrLogin,
	repository.ErrMaxConn:  ErrInternal,
	repository.ErrNotNull:  ErrClientNotNull,
	repository.ErrPassword: ErrClientPassword,
	repository.ErrLogin:    ErrLogin,
}

func ValidationError(err error) error {
	resultErr, ok := arrErr[err]
	if !ok {
		return ErrInternal
	}

	return resultErr
}
