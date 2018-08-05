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
					Datacoding: model.DatacodingPlain,
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
					Datacoding: model.DatacodingPlain,
					UDH:        "050003CC0201",
				},
				{
					Message:    "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB",
					Datacoding: model.DatacodingPlain,
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
					Datacoding: model.DatacodingPlain,
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
					Datacoding: model.DatacodingPlain,
					UDH:        "050003CC0201",
				},
				{
					Message:    "€-----------------",
					Datacoding: model.DatacodingPlain,
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
					Datacoding: model.DatacodingUnicode,
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
					Datacoding: model.DatacodingUnicode,
					UDH:        "050003CC0201",
				},
				{
					Message:    "°°°°",
					Datacoding: model.DatacodingUnicode,
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
					Datacoding: model.DatacodingUnicode,
					UDH:        "050003CC0201",
				},
				{
					Message:    "°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°°",
					Datacoding: model.DatacodingUnicode,
					UDH:        "050003CC0202",
				},
			},
		},
	}

	for _, testCase := range tests {
		splitted := Split(testCase.body)
		if testCase.expectedMessages[0].UDH != "" {
			for i, msg := range testCase.expectedMessages {
				require.Regexp(t, "050003.*", msg.UDH)
				testCase.expectedMessages[i].UDH = ""
			}
			for i := range splitted {
				splitted[i].UDH = ""
			}
		}
		require.Equal(t, testCase.expectedMessages, splitted)
	}
}

func TestGetDataCoding(t *testing.T) {
	tests := []struct {
		body               string
		expectedDatacoding model.Datacoding
	}{
		{
			body:               "",
			expectedDatacoding: model.DatacodingPlain,
		},
		{
			body:               "   Test",
			expectedDatacoding: model.DatacodingPlain,
		},
		{
			body:               "øøøøøøø",
			expectedDatacoding: model.DatacodingPlain,
		},
		{
			body:               "øøøøøøø",
			expectedDatacoding: model.DatacodingPlain,
		},
		{
			body:               "[~]|€",
			expectedDatacoding: model.DatacodingPlain,
		},
		{
			body:               "[~]|€ç",
			expectedDatacoding: model.DatacodingUnicode,
		},
	}
	for _, testCase := range tests {
		require.Equal(t, testCase.expectedDatacoding, getDatacoding(testCase.body))
	}
}
