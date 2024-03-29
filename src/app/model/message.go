package model

import (
	"app/errors"
	"encoding/hex"

	"github.com/nyaruka/phonenumbers"
)

type Datacoding string

const (
	maxOriginatorLength = 11
	maxRecipientLength  = 50

	DatacodingPlain   Datacoding = "plain"
	DatacodingUnicode Datacoding = "unicode"
)

type Message struct {
	Originator string   `json:"originator,omitempty"`
	Body       string   `json:"body,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	UDH        string
	Datacoding Datacoding
}

func (message Message) GetBinaryBody() string {
	src := []byte(message.Body)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
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
