package main

import (
	"app/api"
	"app/model"
	"app/server"
	"log"
	"net/http"

	"github.com/messagebird/go-rest-api"
)

var (
	apiKey string
)

func main() {
	// Create the request channel
	requestChannel := make(chan model.MBSendRequest)
	// Create messagebird client
	if apiKey == "" {
		log.Fatalln("No API Key is provided...")
	}
	client := messagebird.New(apiKey)
	// Create the api
	messagingAPI := api.New(requestChannel, client)
	messagingAPI.StartRequestLoop()

	http.Handle("/sendMessage", server.Method(http.MethodPost, http.HandlerFunc(messagingAPI.SendMessage)))
	log.Print("Starting service...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
