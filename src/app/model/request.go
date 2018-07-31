package model

import (
	"context"

	messagebird "github.com/messagebird/go-rest-api"
)

type MBSendRequest struct {
	ResponseChannel chan *messagebird.Message
	Context         context.Context
	Message         Message
}
