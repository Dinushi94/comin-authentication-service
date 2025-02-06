// internal/common/errors/errors.go
package errors

import "fmt"

type ErrorCode string

const (
	ErrInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrForbidden       ErrorCode = "FORBIDDEN"
	ErrNotFound        ErrorCode = "NOT_FOUND"
	ErrInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDuplicateEntity ErrorCode = "DUPLICATE_ENTITY"
	ErrValidation      ErrorCode = "VALIDATION_ERROR"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewInvalidInputError(msg string, details any) error {
	return AppError{Code: ErrInvalidInput, Message: msg, Details: details}
}

func NewUnauthorizedError(msg string) error {
	return AppError{Code: ErrUnauthorized, Message: msg}
}

func NewForbiddenError(msg string) error {
	return AppError{Code: ErrForbidden, Message: msg}
}

func NewNotFoundError(msg string) error {
	return AppError{Code: ErrNotFound, Message: msg}
}

func NewInternalServerError(msg string) error {
	return AppError{Code: ErrInternalServer, Message: msg}
}

func NewDuplicateEntityError(msg string) error {
	return AppError{Code: ErrDuplicateEntity, Message: msg}
}

func NewValidationError(msg string, details any) error {
	return AppError{Code: ErrValidation, Message: msg, Details: details}
}
