package api

import (
	"app/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/messagebird/go-rest-api"

	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {

	tests := []struct {
		description         string
		request             *model.Message
		messageBirdResponse *messagebird.Message
		expectedBody        string
		expectedStatusCode  int
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
			description: "Happy path",
			request: &model.Message{
				Originator: "Test originator",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			messageBirdResponse: &messagebird.Message{
				Originator: "Test originator",
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
					request.ResponseChannel <- testCase.messageBirdResponse
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
		})
	}
}
