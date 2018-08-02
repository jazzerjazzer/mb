package api

import (
	"app/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/messagebird/go-rest-api"

	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {

	tests := []struct {
		description          string
		request              *model.Message
		messageBirdResponses []*messagebird.Message
		expectedStatusCode   int
		timeout              bool
	}{
		{
			description:        "nil request",
			request:            nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "no request body",
			request:            &model.Message{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "response timeout",
			request: &model.Message{
				Originator: "TestO",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			expectedStatusCode: http.StatusRequestTimeout,
			timeout:            true,
		},
		{
			description: "messagebird api errors",
			request: &model.Message{
				Originator: "TestO",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			messageBirdResponses: []*messagebird.Message{
				&messagebird.Message{
					Errors: []messagebird.Error{
						{Code: http.StatusInternalServerError},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "messagebird nil resp",
			request: &model.Message{
				Originator: "TestO",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			messageBirdResponses: []*messagebird.Message{
				nil,
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "Happy path",
			request: &model.Message{
				Originator: "TestO",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			messageBirdResponses: []*messagebird.Message{
				{
					Originator: "TestO",
					Body:       "Test message",
					Errors:     nil,
					Recipients: messagebird.Recipients{
						TotalCount:     1,
						TotalSentCount: 1,
						Items: []messagebird.Recipient{
							{
								Recipient: 1,
							},
							{
								Recipient: 2,
							},
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "Happy path/unicode/splitted message",
			request: &model.Message{
				Originator: "TestO",
				Body:       "This is a test message. Also it contains some uniçode chars, let's  } make it two messages",
				Recipients: []string{"1", "2"},
			},
			messageBirdResponses: []*messagebird.Message{
				{
					Originator: "TestO",
					Body:       "This is a test message. Also it contains some uniçode chars, let's ",
					Errors:     nil,
					Recipients: messagebird.Recipients{
						TotalCount:     2,
						TotalSentCount: 2,
						Items: []messagebird.Recipient{
							{
								Recipient: 1,
							},
							{
								Recipient: 2,
							},
						},
					},
				},
				{
					Originator: "TestO",
					Body:       " } make it two messages",
					Errors:     nil,
					Recipients: messagebird.Recipients{
						TotalCount:     2,
						TotalSentCount: 2,
						Items: []messagebird.Recipient{
							{
								Recipient: 1,
							},
							{
								Recipient: 2,
							},
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			requestChannel := make(chan model.MBSendRequest, 1)
			messagingAPI := New(requestChannel, nil)
			// Send fake response
			go func() {
				select {
				case request := <-messagingAPI.requests:
					if testCase.timeout {
						// Simulate timeout
						time.Sleep(10 * time.Second)
					} else {
						for _, resp := range testCase.messageBirdResponses {
							request.ResponseChannel <- resp
						}
					}
				}
			}()
			defer close(requestChannel)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(testCase.request)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/sendMessage", &buf)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(messagingAPI.SendMessage)
			handler.ServeHTTP(rr, req)

			require.Equal(t, testCase.expectedStatusCode, rr.Code)

			json.NewEncoder(&buf).Encode(testCase.messageBirdResponses)
			if testCase.expectedStatusCode == http.StatusOK {
				require.Equal(t, buf.String(), rr.Body.String())
			}
		})
	}
}
