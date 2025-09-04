package domainerrors

import "errors"

var (
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrPasswordTooShort   = errors.New("password too short")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenInvalid       = errors.New("token invalid")

	ErrForbidden = errors.New("forbidden")

	ErrConflict = errors.New("conflict")
)
