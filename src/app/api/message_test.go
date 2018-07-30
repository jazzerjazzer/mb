package api

import (
	"app/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage(t *testing.T) {

	tests := []struct {
		description        string
		request            *http.Request
		expectedBody       string
		expectedStatusCode int
	}{
		{
			description: "Happy path",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			requestChannel := make(chan model.MBSendRequest)
			messagingAPI := New(requestChannel, nil)
			defer close(requestChannel)

			rr := httptest.NewRecorder()
			http.HandlerFunc(messagingAPI.SendMessage).ServeHTTP(rr, testCase.request)
		})
	}
}
