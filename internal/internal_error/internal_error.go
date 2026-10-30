package internal_error

import "net/http"

type InternalError struct {
	Status  int
	Message string
}

func NewBadRequestError(message string) *InternalError {
	return &InternalError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) *InternalError {
	return &InternalError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewUnprocessableEntityError(message string) *InternalError {
	return &InternalError{
		Status:  http.StatusUnprocessableEntity,
		Message: message,
	}
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
