package models

type UserRegister struct {
	Name     string `validate:"required,min=5,max=60"`
	Password string `validate:"required,min=8,max=80"`
	Email    string `validate:"required,email"`
}

type UserLogin struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=80"`
}
