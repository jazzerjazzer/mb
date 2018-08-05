package client

import messagebird "github.com/messagebird/go-rest-api"

type Interface interface {
	NewMessage(originator string, recipients []string, body string, msgParams *messagebird.MessageParams) (*messagebird.Message, error)
}
