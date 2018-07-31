package api

import (
	"app/client/mocked"
	"app/model"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
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
				},
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
				},
				Context: ctx,
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
					Body:       "body",
					Recipients: []string{"recipient"},
				},
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
				message.Recipients, message.Body, mock.Anything).Return(testCase.messagebirdResponse, testCase.messagebirdError)
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
