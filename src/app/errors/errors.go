package errors

import "errors"

const (
	errEmptyBody = "Body cannot be empty"
)

// ErrEmptyBody ...
var ErrEmptyBody = errors.New(errEmptyBody)
