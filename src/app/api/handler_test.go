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
				{
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

func TestComposeRequest(t *testing.T) {
	tests := []struct {
		description     string
		givenMessage    model.Message
		expectedRequest []model.MBSendRequest
	}{
		{
			description: "single message/plain",
			givenMessage: model.Message{
				Originator: "TestO",
				Body:       "Test message",
				Recipients: []string{"1", "2"},
			},
			expectedRequest: []model.MBSendRequest{
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       "Test message",
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingPlain,
					},
					MessageType: model.MessageTypeSMS,
				},
			},
		},
		{
			description: "single message/unicode",
			givenMessage: model.Message{
				Originator: "TestO",
				Body:       "Test message: ç",
				Recipients: []string{"1", "2"},
			},
			expectedRequest: []model.MBSendRequest{
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       "Test message: ç",
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingUnicode,
					},
					MessageType: model.MessageTypeSMS,
				},
			},
		},
		{
			description: "multipart message/plain",
			givenMessage: model.Message{
				Originator: "TestO",
				Body:       "This is a long test message. The datacoding is plain, no unicode characters are used. This is a long test message. The datacoding is plain, no unicode characters are used.",
				Recipients: []string{"1", "2"},
			},
			expectedRequest: []model.MBSendRequest{
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       model.Message{Body: "This is a long test message. The datacoding is plain, no unicode characters are used. This is a long test message. The datacoding is plain, no unicode ch"}.GetBinaryBody(),
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingPlain,
					},
					MessageType: model.MessageTypeBinary,
				},
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       model.Message{Body: "aracters are used."}.GetBinaryBody(),
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingPlain,
					},
					MessageType: model.MessageTypeBinary,
				},
			},
		},
		{
			description: "multipart message/unicode",
			givenMessage: model.Message{
				Originator: "TestO",
				Body:       "This is a test message. Also it contains some uniçode chars, let's  } make it two messages",
				Recipients: []string{"1", "2"},
			},
			expectedRequest: []model.MBSendRequest{
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       model.Message{Body: "This is a test message. Also it contains some uniçode chars, let's "}.GetBinaryBody(),
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingPlain,
					},
					MessageType: model.MessageTypeBinary,
				},
				{
					Message: model.Message{
						Originator: "TestO",
						Body:       model.Message{Body: " } make it two messages"}.GetBinaryBody(),
						Recipients: []string{"1", "2"},
						Datacoding: model.DatacodingPlain,
					},
					MessageType: model.MessageTypeBinary,
				},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			composed := composeRequest(testCase.givenMessage, nil, nil)
			if len(testCase.expectedRequest) > 1 {
				for i, comp := range composed {
					require.Regexp(t, "050003.*", comp.Message.UDH)
					composed[i].Message.UDH = ""
				}
			}
			require.Equal(t, testCase.expectedRequest, composed)
		})
	}
}
