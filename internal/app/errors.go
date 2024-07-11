package app

import "fmt"

type AppError struct {
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Err:     err,
	}
}
