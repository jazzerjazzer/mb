package api

import (
	"app/client"
	"app/model"
	"time"
)

type MessageAPI struct {
	requests chan (model.MBSendRequest)
	client   client.Interface
}

func New(requestChannel chan (model.MBSendRequest), c client.Interface) *MessageAPI {
	return &MessageAPI{
		requests: requestChannel,
		client:   c,
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
