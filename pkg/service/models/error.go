package models

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	code    int
	Message string `json:"message"`
}

func NewDoesNotExistError(objectName, property, propertyValue string) *ApiError {
	msg := fmt.Sprintf("%s with %s %s does not exist.", objectName, property, propertyValue)
	return NewApiError(http.StatusNotFound, msg)
}

func NewAlreadyExistsError(objectName, property, propertyValue string) *ApiError {
	msg := fmt.Sprintf("%s with %s %s already exists.", objectName, property, propertyValue)
	return NewApiError(http.StatusConflict, msg)
}

func NewValidationError(message string) *ApiError {
	return NewApiError(http.StatusBadRequest, message)
}

func NewApiError(code int, message string) *ApiError {
	return &ApiError{code: code, Message: message}
}

func (e *ApiError) Code() int {
	return e.code
}

func (e *ApiError) Error() string {
	return e.Message
}
