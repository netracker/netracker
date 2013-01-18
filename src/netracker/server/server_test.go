package server

import (
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"github.com/drewolson/testflight/ws"
	"testing"
)

func withServer(f func(*testflight.Requester)) {
	netrackerServer := New()
	netrackerServer.connectionManager.Run()
	defer netrackerServer.connectionManager.Stop()
	testflight.WithServer(netrackerServer.handler(), f)
}

func TestRoot(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		response := r.Get("/")

		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestWebsocketConnect(t *testing.T) {
	withServer(func(r *testflight.Requester) {
		connection := ws.Connect(r, "/ws")
		message := connection.ReceiveMessage()
		connection.WriteMessage("quit")

		expectedMessage := `{"ActivePlayer":{"Role":"corp"},"InactivePlayer":{"Role":"runner"},"CorpCredits":5,"RunnerCredits":5,"Clicks":0}`

		assert.Equal(t, expectedMessage, message)
	})
}
