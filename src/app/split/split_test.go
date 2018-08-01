package split

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		description      string
		body             string
		expectedMessages []string
	}{
		{
			description:      "plain/160 chars, 1 sms",
			body:             "----------------------------------------------------------------------------------------------------------------------------------------------------------------",
			expectedMessages: []string{"----------------------------------------------------------------------------------------------------------------------------------------------------------------"},
		},
		{
			description: "plain/161 chars, 2 sms",
			body:        "-----------------------------------------------------------------------------------------------------------------------------------------------------------------",
			expectedMessages: []string{
				"---------------------------------------------------------------------------------------------------------------------------------------------------------",
				"--------",
			},
		},
		{
			description: "plain/306 chars, 2 sms",
			body:        "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			expectedMessages: []string{
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
				"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			},
		},
		{
			description: "plain/159 chars, € is the 153th char, 1 sms",
			body:        "--------------------------------------------------------------------------------------------------------------------------------------------------------€------",
			expectedMessages: []string{
				"--------------------------------------------------------------------------------------------------------------------------------------------------------€------",
			},
		},
		{
			description: "plain/170 chars, € is the 153th char, 1 sms",
			body:        "--------------------------------------------------------------------------------------------------------------------------------------------------------€-----------------",
			expectedMessages: []string{
				"--------------------------------------------------------------------------------------------------------------------------------------------------------",
				"€-----------------",
			},
		},
		{
			description: "unicode/70 chars, 1 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []string{
				"°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			},
		},
		{
			description: "unicode/71 chars, 2 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []string{
				"°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
				"°°°°",
			},
		},
		{
			description: "unicode/134 chars, 2 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []string{
				"°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
				"°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			},
		},
	}

	for _, testCase := range tests {
		require.Equal(t, testCase.expectedMessages, Split(testCase.body))
	}
}

func TestGetDataCoding(t *testing.T) {
	tests := []struct {
		body               string
		expectedDatacoding Datacoding
	}{
		{
			body:               "",
			expectedDatacoding: datacodingPlain,
		},
		{
			body:               "   Test",
			expectedDatacoding: datacodingPlain,
		},
		{
			body:               "øøøøøøø",
			expectedDatacoding: datacodingPlain,
		},
		{
			body:               "øøøøøøø",
			expectedDatacoding: datacodingPlain,
		},
		{
			body:               "[~]|€",
			expectedDatacoding: datacodingPlain,
		},
		{
			body:               "[~]|€ç",
			expectedDatacoding: datacodingUnicode,
		},
	}
	for _, testCase := range tests {
		require.Equal(t, testCase.expectedDatacoding, getDatacoding(testCase.body))
	}
}
