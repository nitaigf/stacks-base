package services

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(statusCode int, code string, message string, err error) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message, Err: err}
}

func AsAppError(err error) (*AppError, bool) {
	var appError *AppError
	if errors.As(err, &appError) {
		return appError, true
	}

	return nil, false
}

func InternalError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, "internal_error", message, err)
}