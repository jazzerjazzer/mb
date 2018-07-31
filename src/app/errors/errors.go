package errors

import "errors"

const (
	errNilMessage                  = "Message is nil"
	errEmptyBody                   = "Body cannot be empty"
	errEmptyOriginator             = "Originator cannot not be empty"
	errMaxOriginatorLenghtExceeded = "Max originator length exceeded"
	errEmptyRecipients             = "Recipients cannot be empty"
	errMaxRecipientsExceeded       = "Max recipients exceeded"
)

// ErrNilMessage ...
var ErrNilMessage = errors.New(errNilMessage)

// ErrEmptyBody ...
var ErrEmptyBody = errors.New(errEmptyBody)

// ErrMaxOriginatorLenghtExceeded ...
var ErrMaxOriginatorLenghtExceeded = errors.New(errMaxOriginatorLenghtExceeded)

// ErrEmptyOriginator ...
var ErrEmptyOriginator = errors.New(errEmptyOriginator)

// ErrMaxRecipientsExceeded ...
var ErrMaxRecipientsExceeded = errors.New(errMaxRecipientsExceeded)

// ErrEmptyRecipients ...
var ErrEmptyRecipients = errors.New(errEmptyRecipients)
