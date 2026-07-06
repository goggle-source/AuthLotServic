package grpc

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/goggle-source/authLotServic/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Err struct {
	Code codes.Code
	Err  error
}

func ValidationError(err error) error {
	if errors.Is(err, domain.ErrEmail) {
		return status.Error(codes.InvalidArgument, "such an email has already been registered")
	}
	if errors.Is(err, domain.ErrGenerateJWT) {
		return status.Error(codes.Internal, "internal error")
	}
	if errors.Is(err, domain.ErrPasswordOrEmail) {
		return status.Error(codes.InvalidArgument, "err in email or password")
	}
	if errors.Is(err, domain.ErrUserNoFound) {
		return status.Error(codes.InvalidArgument, "err in email or password")
	}

	return status.Error(codes.Internal, "internal error")
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
