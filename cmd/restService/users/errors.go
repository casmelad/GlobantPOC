package users

import "errors"

var (
	ErrNotFound          error = errors.New("not found")
	ErrInternalFailure   error = errors.New("bad request")
	ErrInvalidInput      error = errors.New("invalid data")
	ErrUserAlreadyExists error = errors.New("already exists")
)

type AppError struct {
	errorMessage error
}

func WrapError(msg error) AppError {
	return AppError{
		errorMessage: msg,
	}
}

func (a AppError) error() error {
	return a.errorMessage
}
