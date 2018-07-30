package model

import (
	"context"

	messagebird "github.com/messagebird/go-rest-api"
)

type Message struct {
	Originator string   `json:"originator,omitempty"`
	Body       string   `json:"body,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
}

func (m *Message) Validate() error {
	return nil
}

type MBSendRequest struct {
	ResponseChannel chan *messagebird.Message
	Context         context.Context
	Message         Message
}
