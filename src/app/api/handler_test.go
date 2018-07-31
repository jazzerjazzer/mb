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
		description         string
		request             *model.Message
		messageBirdResponse *messagebird.Message
		expectedStatusCode  int
		timeout             bool
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
			messageBirdResponse: &messagebird.Message{
				Errors: []messagebird.Error{
					{Code: http.StatusInternalServerError},
				},
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
			messageBirdResponse: &messagebird.Message{
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
						time.Sleep(5 * time.Second)
					} else {
						request.ResponseChannel <- testCase.messageBirdResponse
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

			json.NewEncoder(&buf).Encode(testCase.messageBirdResponse)
			if testCase.expectedStatusCode == http.StatusOK {
				require.Equal(t, buf.String(), rr.Body.String())
			}
		})
	}
}
