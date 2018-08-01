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
	var message model.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	for i := 0; i < 2; i++ {
		select {
		case r := <-response:
			composeResponse(w, r)
		case <-time.After(10 * time.Second):
			cancel()
			w.WriteHeader(http.StatusRequestTimeout)
		}
	}
}

func composeResponse(w http.ResponseWriter, r *messagebird.Message) {
	if r == nil || len(r.Errors) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
