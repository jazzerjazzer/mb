package model

import (
	"context"

	messagebird "github.com/messagebird/go-rest-api"
)

type MessageType string

const (
	MessageTypeSMS    MessageType = "sms"
	MessageTypeBinary MessageType = "binary"
)

type MBSendRequest struct {
	ResponseChannel chan *messagebird.Message
	Context         context.Context
	Message         Message
	MessageType     MessageType
}
