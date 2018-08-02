package split

import (
	"app/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		description      string
		body             string
		expectedMessages []model.Split
	}{
		{
			description: "plain/160 chars, 1 sms",
			body:        "----------------------------------------------------------------------------------------------------------------------------------------------------------------",
			expectedMessages: []model.Split{
				{
					Message:    "----------------------------------------------------------------------------------------------------------------------------------------------------------------",
					Datacoding: datacodingPlain,
					UDH:        "",
				},
			},
		},
		{
			description: "plain/306 chars, 2 sms",
			body:        "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
			expectedMessages: []model.Split{
				{
					Message:    "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
					Datacoding: datacodingPlain,
					UDH:        "050003CC0201",
				},
				{
					Message:    "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
					Datacoding: datacodingPlain,
					UDH:        "050003CC0202",
				},
			},
		},
		{
			description: "plain/159 chars, € is the 153th char, 1 sms",
			body:        "--------------------------------------------------------------------------------------------------------------------------------------------------------€------",
			expectedMessages: []model.Split{
				{
					Message:    "--------------------------------------------------------------------------------------------------------------------------------------------------------€------",
					Datacoding: datacodingPlain,
					UDH:        "",
				},
			},
		},
		{
			description: "plain/170 chars, € is the 153th char, 2 sms",
			body:        "--------------------------------------------------------------------------------------------------------------------------------------------------------€-----------------",
			expectedMessages: []model.Split{
				{
					Message:    "--------------------------------------------------------------------------------------------------------------------------------------------------------",
					Datacoding: datacodingPlain,
					UDH:        "050003CC0201",
				},
				{
					Message:    "€-----------------",
					Datacoding: datacodingPlain,
					UDH:        "050003CC0202",
				},
			},
		},
		{
			description: "unicode/70 chars, 1 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []model.Split{
				{
					Message:    "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
					Datacoding: datacodingUnicode,
					UDH:        "",
				},
			},
		},
		{
			description: "unicode/71 chars, 2 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []model.Split{
				{
					Message:    "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
					Datacoding: datacodingUnicode,
					UDH:        "050003CC0201",
				},
				{
					Message:    "°°°°",
					Datacoding: datacodingUnicode,
					UDH:        "050003CC0202",
				},
			},
		},
		{
			description: "unicode/71 chars, 2 sms",
			body:        "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
			expectedMessages: []model.Split{
				{
					Message:    "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
					Datacoding: datacodingUnicode,
					UDH:        "050003CC0201",
				},
				{
					Message:    "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
					Datacoding: datacodingUnicode,
					UDH:        "050003CC0202",
				},
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
