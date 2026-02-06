package models

type UserAddDatabase struct {
	Name         string
	PasswordHash []byte
	Email        string
	Id           string
}

type UserValidateInDatabase struct {
	Email    string
	Password string
}
