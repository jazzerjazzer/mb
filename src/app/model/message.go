package model

import (
	"app/errors"

	"github.com/nyaruka/phonenumbers"
)

const (
	maxOriginatorLength = 11
	maxRecipientLength  = 50
)

type Message struct {
	Originator string   `json:"originator,omitempty"`
	Body       string   `json:"body,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
}

func (m *Message) Validate() error {
	if m == nil {
		return errors.ErrNilMessage
	}

	if m.Originator == "" {
		return errors.ErrEmptyOriginator
	}

	_, err := phonenumbers.Parse(m.Originator, "")
	if err != nil {
		// Fallback to alphanumeric string
		if len(m.Originator) > maxOriginatorLength {
			return errors.ErrMaxOriginatorLenghtExceeded
		}
	}

	if m.Body == "" {
		return errors.ErrEmptyBody
	}

	if len(m.Recipients) == 0 {
		return errors.ErrEmptyRecipients
	}

	if len(m.Recipients) > maxRecipientLength {
		return errors.ErrMaxRecipientsExceeded
	}

	return nil
}
