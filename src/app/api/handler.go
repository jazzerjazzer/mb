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

	messages := split.Split(message.Body)
	for _, msg := range messages {
		splitted := model.Message{
			Recipients: message.Recipients,
			Originator: message.Originator,
			Body:       msg.Message,
			UDH:        msg.UDH,
			Datacoding: msg.Datacoding,
		}
		request := model.MBSendRequest{
			ResponseChannel: response,
			Context:         ctx,
			Message:         splitted,
		}
		api.requests <- request
	}

	var responses []*messagebird.Message
	// Wait for all responses
	for i := 0; i < len(messages); i++ {
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
