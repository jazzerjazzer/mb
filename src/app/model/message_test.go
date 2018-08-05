package model

import (
	"app/errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestGetBinaryBody(t *testing.T) {
	tests := []struct {
		description        string
		givenMessage       Message
		expectedBinaryBody string
	}{
		{
			description: "empty base case",
		},
		{
			description: "Base7GSMBody",
			givenMessage: Message{
				Body: "This is a Base7GSM body",
			},
			expectedBinaryBody: "54686973206973206120426173653747534d20626f6479",
		},
		{
			description: "BaseAndExtended7GSMBody",
			givenMessage: Message{
				Body: "This is a Base7GSM body with {{||extended characters||}}",
			},
			expectedBinaryBody: "54686973206973206120426173653747534d20626f64792077697468207b7b7c7c657874656e64656420636861726163746572737c7c7d7d",
		},
		{
			description: "UnicodeBody",
			givenMessage: Message{
				Body: "This is a Base7GSM body with unicode characters like ç",
			},
			expectedBinaryBody: "54686973206973206120426173653747534d20626f6479207769746820756e69636f64652063686172616374657273206c696b6520c3a7",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.expectedBinaryBody, testCase.givenMessage.GetBinaryBody())
		})
	}
}
