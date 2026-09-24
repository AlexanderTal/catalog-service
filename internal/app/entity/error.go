package entity

import "net/http"

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string   { return e.Message }
func (e *AppError) HTTPStatus() int { return e.Status }

func NewAppError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

var (
	ErrNotFound            = NewAppError(http.StatusNotFound, "not found")
	ErrAlreadyExists       = NewAppError(http.StatusConflict, "already exists")
	ErrCategoryHasProducts = NewAppError(http.StatusConflict, "category has linked products")
	ErrIncorrectParameters = NewAppError(http.StatusBadRequest, "incorrect parameters")
)
