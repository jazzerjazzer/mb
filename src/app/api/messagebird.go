package api

import (
	"app/model"
	"log"

	"github.com/messagebird/go-rest-api"
)

const (
	MessageTypeBinary = "binary"
	TypeDetailUDH     = "udh"
)

func (mb *MessageAPI) Send(req model.MBSendRequest) {
	params := &messagebird.MessageParams{
		Type:        MessageTypeBinary,
		TypeDetails: messagebird.TypeDetails{TypeDetailUDH: req.Message.UDH},
		DataCoding:  req.Message.Datacoding,
	}

	resp, err := mb.client.NewMessage(req.Message.Originator, req.Message.Recipients, req.Message.GetBinaryBody(), params)
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
