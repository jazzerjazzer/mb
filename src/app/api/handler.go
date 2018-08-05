package api

import (
	"app/model"
	"app/split"
	"context"
	"encoding/json"
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// SendMessage is an endpoint handler
// which constructs and sends request to requests channel for the message to be sent to Messagebird backend.
func (api *MessageAPI) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the request
	var message model.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate the request
	if err := message.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make(chan *messagebird.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer close(response)

	// Compose and send request to the requests channel
	requests := composeRequest(message, response, ctx)
	for _, request := range requests {
		api.requests <- request
	}

	var responses []*messagebird.Message
	// Wait for all responses
	for i := 0; i < len(requests); i++ {
		select {
		case r := <-response:
			responses = append(responses, r)
		case <-time.After(10 * time.Second):
			cancel()
			w.WriteHeader(http.StatusRequestTimeout)
		}
	}
	composeResponse(w, responses)
}

func composeRequest(message model.Message, response chan *messagebird.Message, ctx context.Context) []model.MBSendRequest {
	// Split the messages
	messages := split.Split(message.Body)
	// If single message, send as an sms, not as binary
	if len(messages) == 1 {
		message.Datacoding = messages[0].Datacoding
		request := model.MBSendRequest{
			ResponseChannel: response,
			Context:         ctx,
			Message:         message,
			MessageType:     model.MessageTypeSMS,
		}
		return []model.MBSendRequest{request}
	}
	var requests []model.MBSendRequest
	for _, msg := range messages {
		splitted := model.Message{
			Recipients: message.Recipients,
			Originator: message.Originator,
			Body:       msg.Message,
			UDH:        msg.UDH,
			Datacoding: model.DatacodingPlain,
		}
		// Get the binary body instead of raw
		splitted.Body = splitted.GetBinaryBody()
		request := model.MBSendRequest{
			ResponseChannel: response,
			Context:         ctx,
			Message:         splitted,
			MessageType:     model.MessageTypeBinary,
		}
		requests = append(requests, request)
	}
	return requests
}

func composeResponse(w http.ResponseWriter, responses []*messagebird.Message) {
	for _, r := range responses {
		if r == nil {
			// Unexpected response from messagebird
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
