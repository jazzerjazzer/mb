package api

import (
	"app/model"
	"time"

	"github.com/messagebird/go-rest-api"
)

type MessageAPI struct {
	requests chan (model.MBSendRequest)
	client   *messagebird.Client
}

func New(requestChannel chan (model.MBSendRequest), client *messagebird.Client) *MessageAPI {
	return &MessageAPI{
		requests: requestChannel,
		client:   client,
	}
}

func (m *MessageAPI) StartRequestLoop() {
	throttle := time.Tick(time.Second)
	go func() {
		for req := range m.requests {
			<-throttle
			go m.Send(req)
		}
	}()
}
