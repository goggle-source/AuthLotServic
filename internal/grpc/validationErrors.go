package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/goggle-source/authLotServic/internal/servic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Err struct {
	Code codes.Code
	Err  error
}

var arrErr = map[error]Err{
	servic.ErrGenerateJWT:    {Code: codes.Internal, Err: ErrInternal},
	servic.ErrClientEnuqie:   {Code: codes.InvalidArgument, Err: ErrClientEnuqie},
	servic.ErrClientNotNull:  {Code: codes.InvalidArgument, Err: ErrClientNotNull},
	servic.ErrValidateToken:  {Code: codes.InvalidArgument, Err: ErrNoValidToken},
	servic.ErrClientPassword: {Code: codes.InvalidArgument, Err: ErrClientPassword},
	servic.ErrNoFound:        {Code: codes.InvalidArgument, Err: ErrNoFound},
}

func ValidationError(err error) error {
	result, ok := arrErr[err]
	if !ok {
		return status.Error(codes.Internal, ErrInternal.Error())
	}

	return status.Error(result.Code, result.Err.Error())
}

func ValidationErrValidator(err error) error {
	var errResult string
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errResult = errResult + " requred:" + fieldError.Error()
			case "email":
				errResult = errResult + " email is invalid:" + fieldError.Error()
			case "min":
				errResult = errResult + " error min:" + fieldError.Error()
			case "max":
				errResult = errResult + " error max:" + fieldError.Error()
			}
		}
	}

	return status.Error(codes.InvalidArgument, errResult)
}
