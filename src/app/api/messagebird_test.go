package api

import (
	"app/client/mocked"
	"app/model"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/messagebird/go-rest-api"
)

func TestSend(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tests := []struct {
		description         string
		request             model.MBSendRequest
		messagebirdResponse *messagebird.Message
		messagebirdError    error
		contextCancel       context.CancelFunc
		isCancelled         bool
	}{
		{
			description: "messagebird api call error",
			request: model.MBSendRequest{
				ResponseChannel: make(chan *messagebird.Message),
				Context:         context.Background(),
				Message: model.Message{
					Originator: "originator",
					Body:       "body",
					Recipients: []string{"recipient"},
					Datacoding: model.DatacodingPlain,
				},
				MessageType: model.MessageTypeSMS,
			},
			messagebirdResponse: nil,
			messagebirdError:    errors.New("API error"),
		},
		{
			description: "context cancelled",
			request: model.MBSendRequest{
				ResponseChannel: make(chan *messagebird.Message),
				Message: model.Message{
					Originator: "originator",
					Body:       "body",
					Recipients: []string{"recipient"},
					Datacoding: model.DatacodingPlain,
				},
				MessageType: model.MessageTypeSMS,
				Context:     ctx,
			},
			contextCancel: cancel,
			isCancelled:   true,
		},
		{
			description: "messagebird api call success",
			request: model.MBSendRequest{
				ResponseChannel: make(chan *messagebird.Message),
				Context:         context.Background(),
				Message: model.Message{
					Originator: "originator",
					Body:       "body ç",
					Recipients: []string{"recipient"},
					UDH:        "TestUDH",
					Datacoding: model.DatacodingUnicode,
				},
				MessageType: model.MessageTypeSMS,
			},
			messagebirdResponse: &messagebird.Message{
				Originator: "originator",
				Body:       "body",
			},
			messagebirdError: nil,
		},
		{
			description: "messagebird api call success/binary",
			request: model.MBSendRequest{
				ResponseChannel: make(chan *messagebird.Message),
				Context:         context.Background(),
				Message: model.Message{
					Originator: "originator",
					Body:       "body",
					Recipients: []string{"recipient"},
					UDH:        "TestUDH",
					Datacoding: model.DatacodingPlain,
				},
				MessageType: model.MessageTypeBinary,
			},
			messagebirdResponse: &messagebird.Message{
				Originator: "originator",
				Body:       "body",
			},
			messagebirdError: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			defer close(testCase.request.ResponseChannel)
			message := testCase.request.Message

			// Mock messagebird calls
			mockedClient := new(mocked.Interface)
			mockedClient.On("NewMessage", message.Originator,
				message.Recipients, message.Body, &messagebird.MessageParams{
					Type:        string(testCase.request.MessageType),
					TypeDetails: messagebird.TypeDetails{TypeDetailUDH: message.UDH},
					DataCoding:  string(message.Datacoding),
				}).Return(testCase.messagebirdResponse, testCase.messagebirdError)
			messagingAPI := New(nil, mockedClient)

			// Test cancelled context path
			if testCase.isCancelled {
				testCase.contextCancel()
				require.NotPanics(t, func() {
					messagingAPI.Send(testCase.request)
				})
				return
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func(w *sync.WaitGroup) {
				resp := <-testCase.request.ResponseChannel
				require.Equal(t, testCase.messagebirdResponse, resp)
				w.Done()
			}(&wg)

			messagingAPI.Send(testCase.request)
			wg.Wait()
		})
	}
}
