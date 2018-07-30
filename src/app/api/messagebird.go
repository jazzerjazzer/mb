package api

import (
	"app/model"
	"log"
)

func (mb *MessageAPI) Send(req model.MBSendRequest) {
	resp, err := mb.client.NewMessage(req.Message.Originator, req.Message.Recipients, req.Message.Body, nil)
	if err != nil {
		log.Printf("Cannot send message to MessageBird backend: %+v -- %v", err, resp.Errors)
		resp = nil
	}
	select {
	case <-req.Context.Done():
		return
	default:
		req.ResponseChannel <- resp
	}
}
