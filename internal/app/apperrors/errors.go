package apperrors

import "errors"

type validationError struct{ message string }
type notFoundError struct{ message string }
type conflictError struct{ message string }

func (e validationError) Error() string { return e.message }
func (e notFoundError) Error() string   { return e.message }
func (e conflictError) Error() string   { return e.message }

func NewValidation(message string) error { return validationError{message: message} }
func NewNotFound(message string) error   { return notFoundError{message: message} }
func NewConflict(message string) error   { return conflictError{message: message} }

func IsValidation(err error) bool {
	var target validationError
	return errors.As(err, &target)
}

func IsNotFound(err error) bool {
	var target notFoundError
	return errors.As(err, &target)
}

func IsConflict(err error) bool {
	var target conflictError
	return errors.As(err, &target)
}
