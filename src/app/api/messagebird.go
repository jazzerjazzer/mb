package api

import (
	"app/model"
	"log"

	"github.com/messagebird/go-rest-api"
)

const (
	TypeDetailUDH = "udh"
)

// Send constructs the message params, sends the message to Messagebird backend and
// returns the API response to the respective response channel.
func (mb *MessageAPI) Send(req model.MBSendRequest) {
	params := &messagebird.MessageParams{
		Type:        string(req.MessageType),
		TypeDetails: messagebird.TypeDetails{TypeDetailUDH: req.Message.UDH},
		DataCoding:  string(req.Message.Datacoding),
	}

	resp, err := mb.client.NewMessage(req.Message.Originator, req.Message.Recipients, req.Message.Body, params)
	if err != nil {
		log.Printf("Cannot send message to MessageBird backend: %+v", err)
		resp = nil
	}
	select {
	case <-req.Context.Done():
		return
	default:
		req.ResponseChannel <- resp
	}
}
