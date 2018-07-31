package model

import (
	"app/errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO: Test ASCII and Unicode chars and their lengths

func TestValidate(t *testing.T) {
	var maxRecipients []string
	for i := 0; i < 51; i++ {
		maxRecipients = append(maxRecipients, strconv.Itoa(i))
	}
	tests := []struct {
		description   string
		message       *Message
		expectedError error
	}{
		{
			description:   "nil message",
			message:       nil,
			expectedError: errors.ErrNilMessage,
		},
		{
			description:   "empty originator",
			message:       &Message{},
			expectedError: errors.ErrEmptyOriginator,
		},
		{
			description: "invalid phone number, exceeded alphanum length",
			message: &Message{
				Originator: "TESTORIGINATOR+TESTORIGINATOR",
			},
			expectedError: errors.ErrMaxOriginatorLenghtExceeded,
		},
		{
			description: "valid phone number, empty body",
			message: &Message{
				Originator: "+4915166962555",
			},
			expectedError: errors.ErrEmptyBody,
		},
		{
			description: "empty body",
			message: &Message{
				Originator: "+4915166962555",
			},
			expectedError: errors.ErrEmptyBody,
		},
		{
			description: "empty recipients",
			message: &Message{
				Originator: "+4915166962555",
				Body:       "Test body",
			},
			expectedError: errors.ErrEmptyRecipients,
		},
		{
			description: "more than allowed recipients",
			message: &Message{
				Originator: "+4915166962555",
				Body:       "Test body",
				Recipients: maxRecipients,
			},
			expectedError: errors.ErrMaxRecipientsExceeded,
		},
		{
			description: "valid model",
			message: &Message{
				Originator: "+4915166962555",
				Body:       "Test body",
				Recipients: []string{"recipient1"},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.expectedError, testCase.message.Validate())
		})
	}
}
