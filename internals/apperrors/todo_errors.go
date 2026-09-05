package apperrors

import "errors"

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrInvalidID      = errors.New("invalid task id")
	ErrInvalidRequest = errors.New("invalid request body")
	ErrTitleEmpty     = errors.New("title should not be empty string")
	ErrInvalidStatus  = errors.New("invalid status requested")
	ErrTitleTooLong   = errors.New("title is too long, 200 chars max")
)
