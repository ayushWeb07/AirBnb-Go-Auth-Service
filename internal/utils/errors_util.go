package utils

import "net/http"

type AppError struct {
	Message    string
	StatusCode int
	Success    bool
}

func (appError *AppError) Error() string {
	return appError.Message
}

func NotFound(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func InternalServerError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}
