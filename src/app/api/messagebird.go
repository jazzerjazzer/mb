package api

import (
	"app/model"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/messagebird/go-rest-api"
)

func (mb *MessageAPI) Send(req model.MBSendRequest) {
	params := &messagebird.MessageParams{Type: "binary", TypeDetails: messagebird.TypeDetails{"udh": req.Message.UDH},
		DataCoding: req.Message.Datacoding}

	src := []byte(req.Message.Body)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	fmt.Printf("UDH: %+v\n\n", req.Message.UDH)
	resp, err := mb.client.NewMessage(req.Message.Originator, req.Message.Recipients, string(dst), params)
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
