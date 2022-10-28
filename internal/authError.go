package internal

import "errors"

var (
	CredentialsInUseError = errors.New("username already in use")
	WrongCredentials      = errors.New("no such user or wrong password")
)
